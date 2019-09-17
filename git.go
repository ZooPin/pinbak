package pinbak

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"time"
)

type GitHelper struct {
	Path   string
	Config *Config
}

func CreateGit(config Config) GitHelper {
	return GitHelper{
		Path:   config.path,
		Config: &config,
	}
}

func (g GitHelper) Clone(name string, url string) error {
	var homePath, _ = os.UserHomeDir()
	homePath = fmt.Sprint(homePath, "/.pinbak/", name)
	_, err := git.PlainClone(homePath, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          os.Stdout,
	})

	if err != nil {
		return err
	}

	g.Config.AddRepository(name, url)

	return nil
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
