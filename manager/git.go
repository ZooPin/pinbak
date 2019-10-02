package manager

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"time"
)

type GitHelper struct {
	Path   string
	Config Config
}

func CreateGit(config Config) GitHelper {
	return GitHelper{
		Path:   config.path,
		Config: config,
	}
}

func (g GitHelper) Clone(name string, url string) error {
	homePath := fmt.Sprint(g.Path, "/", name)
	_, err := git.PlainClone(homePath, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return err
	}

	err = g.Config.AddRepository(name, url)

	return err
}

func (g GitHelper) createPath(repoName string) string {
	return fmt.Sprint(g.Path, "/", repoName)
}

func (g GitHelper) CommitAndPush(repoName string) error {
	err := g.Commit(repoName)
	if err != nil {
		return err
	}

	err = g.Push(repoName)
	return err
}

func (g GitHelper) Pull(repoName string) error {
	r, err := git.PlainOpen(g.createPath(repoName))
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
	})

	if err == git.NoErrAlreadyUpToDate {
		err = nil
	}
	return err
}

func (g GitHelper) Push(repoName string) error {
	r, err := git.PlainOpen(g.createPath(repoName))
	if err != nil {
		return err
	}

	err = r.Push(&git.PushOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (g GitHelper) Commit(repoName string) error {
	r, err := git.PlainOpen(g.createPath(repoName))
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	_, err = w.Commit("Backup.", &git.CommitOptions{
		Author: &object.Signature{
			Name:  g.Config.Name,
			Email: g.Config.Email,
			When:  time.Now(),
		},
	})
	return err
}

func (g GitHelper) Remove(repoName string, path string) error {
	r, err := git.PlainOpen(g.createPath(repoName))
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Remove(path)
	if err != nil {
		return err
	}

	return nil
}
