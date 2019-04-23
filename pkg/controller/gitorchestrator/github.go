package gitorchestrator

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

type Github struct {
	GitBase
	ctx        context.Context
	Client     *github.Client
	Visibility bool
}

func (g Github) setEndpoint(endpoint string) error {
	// stub
	return nil
}

func (g Github) setNamespace() error {
	// stub
	return nil
}

func (g Github) setVisibility() error {
	if g.GitBase.Visibility == "private" {
		g.Visibility = true
		return nil
	} else if g.GitBase.Visibility == "public" {
		g.Visibility = false
		return nil
	} else {
		return fmt.Errorf("missing/wrong visibility level")
	}
}

func (g Github) doesProjectExist() bool {
	project, _, _ := g.Client.Repositories.Get(g.ctx, g.GitBase.Namespace, g.GitBase.RepositoryName)
	if project == nil {
		return false
	} else {
		return true
	}
}

func (g Github) createProject() error {
	repo := &github.Repository{
		Name:        github.String(g.GitBase.RepositoryName),
		Description: &g.GitBase.Description,
		Private:     github.Bool(g.Visibility),
	}

	_, _, err := g.Client.Repositories.Create(g.ctx, g.GitBase.Organisation, repo)

	return err
}
