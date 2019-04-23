package gitorchestrator

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	gitorchestratorv1alpha1 "github.com/wiardvanrij/gitorchestrator/pkg/apis/gitorchestrator/v1alpha1"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_gitorchestrator")

// Add creates a new GitOrchestrator Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileGitOrchestrator{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("gitorchestrator-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource GitOrchestrator
	err = c.Watch(&source.Kind{Type: &gitorchestratorv1alpha1.GitOrchestrator{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &gitorchestratorv1alpha1.GitOrchestrator{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileGitOrchestrator{}

// ReconcileGitOrchestrator reconciles a GitOrchestrator object
type ReconcileGitOrchestrator struct {
	client client.Client
	scheme *runtime.Scheme
}

// Interface for the gitclients/functions
type gitClient interface {
	setEndpoint(string) error
	setNamespace() error
	setVisibility() error
	doesProjectExist() bool
	createProject() error
}

// Base struct for the gitClients
type GitBase struct {
	Namespace      string
	Organisation   string
	RepositoryName string
	Description    string
	Visibility     string
}

func (r *ReconcileGitOrchestrator) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	// Start of basic inits & reconcile
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	infoLogger := log.WithName("DEBUG")
	reqLogger.Info("Reconciling GitOrchestrator")

	// Fetch the GitOrchestrator instance
	instance := &gitorchestratorv1alpha1.GitOrchestrator{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Retreive the token to be used for the API calls
	token, set := os.LookupEnv("SECRET_TOKEN")
	if !set {
		reqLogger.Error(fmt.Errorf("Missing access token"), "set SECRET_TOKEN")
		return reconcile.Result{}, nil
	}

	endpoint, _ := os.LookupEnv("END_POINT")

	// Initial check if we did not already processed this repository
	for _, s := range instance.Status.Repositories {
		// This is required because repo titles are not case-sensitive, it's underlaying path is...
		if strings.ToLower(s) == strings.ToLower(instance.Spec.RepositoryName) {
			infoLogger.Info("Not processing repository " + s + ", I have created this before.")
			return reconcile.Result{}, nil
		}
	}

	var gb = GitBase{
		RepositoryName: instance.Spec.RepositoryName,
		Namespace:      instance.Spec.RepositoryNamespace,
		Description:    instance.Spec.Description,
		Organisation:   instance.Spec.Organisation,
		Visibility:     instance.Spec.Visibility,
	}

	var gc gitClient

	// Init of various implementations should be refactored.
	if instance.Spec.GitType == "github" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)

		gc = Github{
			GitBase: gb,
			ctx:     ctx,
			Client:  github.NewClient(tc),
		}
	} else if instance.Spec.GitType == "gitlab" {
		gc = Gitlab{
			GitBase: gb,
			Client:  gitlab.NewClient(nil, token),
		}
	}

	endPointErr := gc.setEndpoint(endpoint)
	if endPointErr != nil {
		reqLogger.Error(endPointErr, "Endpoint error")
		return reconcile.Result{}, nil
	}

	visibilityErr := gc.setVisibility()
	if visibilityErr != nil {
		reqLogger.Error(visibilityErr, "Visibility error: check for a valid visibility level")
		return reconcile.Result{}, nil
	}

	// Truth; this is a bit odd. We have the namespace NAME but gitlab requires an ID.
	// Therefor we require an API call to set this object for access to it's ID.
	namespaceErr := gc.setNamespace()
	if namespaceErr != nil {
		reqLogger.Error(namespaceErr, "Namespace error: check if the repository namespace exist")
		return reconcile.Result{}, nil
	}

	// Check if the project exist, if not; create it.
	if !gc.doesProjectExist() {
		createError := gc.createProject()
		if createError != nil {
			reqLogger.Error(createError, "Failure on creating project")
			return reconcile.Result{}, nil
		}
	}

	// Update the repository to our objects status so we do not keep on trying to create it
	instance.Status.Repositories = append(instance.Status.Repositories, instance.Spec.RepositoryName)
	updateErr := r.client.Status().Update(context.TODO(), instance)
	if updateErr != nil {
		reqLogger.Error(updateErr, "failure on adding the repository to my status")
	}

	return reconcile.Result{}, nil

}
