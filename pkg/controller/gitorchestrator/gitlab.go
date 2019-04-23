package gitorchestrator

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type Gitlab struct {
	GitBase
	Client *gitlab.Client
}

var Namespace *gitlab.Namespace

func (g Gitlab) setEndpoint(endpoint string) error {
	if endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}

	return g.Client.SetBaseURL(endpoint)
}

func (g Gitlab) setVisibility() error {

	if g.GitBase.Visibility == "private" || g.GitBase.Visibility == "internal" || g.GitBase.Visibility == "public" {
		return nil
	} else {
		return fmt.Errorf("missing/wrong visibility level")
	}
}

// So we have to search for our namespace. In obvious cases it should just return 1 result.
// Unless we require the user to provide an actual ID, we must use this for the creation of the repository eventually.
func (g Gitlab) setNamespace() error {

	opt := &gitlab.ListNamespacesOptions{Search: gitlab.String(g.Namespace)}
	namespaces, _, _ := g.Client.Namespaces.ListNamespaces(opt)

	count := len(namespaces)
	if count == 1 {
		Namespace = namespaces[0]
		return nil
	} else if count == 0 {
		return fmt.Errorf("Namespace not found")
	} else {
		return fmt.Errorf("Found to many namespaces..")
	}
}

// Try to fetch the project, if we got it, it does exist
func (g Gitlab) doesProjectExist() bool {
	projectOpt := &gitlab.GetProjectOptions{}
	projectPath := Namespace.Name + "/" + g.RepositoryName
	project, _, _ := g.Client.Projects.GetProject(projectPath, projectOpt)

	if project == nil {
		return false
	} else {
		return true
	}
}

func (g Gitlab) createProject() error {

	p := &gitlab.CreateProjectOptions{
		Name:                 gitlab.String(g.RepositoryName),
		Description:          gitlab.String(g.Description),
		MergeRequestsEnabled: gitlab.Bool(true),
		SnippetsEnabled:      gitlab.Bool(true),
		NamespaceID:          &Namespace.ID,
		Visibility:           gitlab.Visibility(gitlab.VisibilityValue(g.Visibility)),
	}
	_, _, err := g.Client.Projects.CreateProject(p)
	return err
}
