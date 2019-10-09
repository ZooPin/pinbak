package manager

import (
	"errors"
	"github.com/otiai10/copy"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Mover struct {
	Config Config
	Git    GitHelper
	Path   string
	Index  map[string]Index
}

func CreateMover(config Config, git GitHelper) Mover {
	return Mover{
		Config: config,
		Git:    git,
		Path:   config.path,
	}
}

const moverHome = "{HOME}"

func (m *Mover) checkIndex(repoName string) (Index, error) {
	if m.Index == nil {
		m.Index = make(map[string]Index)
	}
	index, ok := m.Index[repoName]
	if !ok {
		index, err := openIndex(m.Path, repoName)
		if err != nil {
			return index, err
		}
		m.Index[repoName] = index
		return index, nil
	}
	return index, nil
}

func (m Mover) Add(path string, repoName string) error {
	err := m.Git.Pull(repoName)
	if err != nil {
		return err
	}

	absPath, err := m.absolutePath(path)
	if err != nil {
		return err
	}

	index, err := m.checkIndex(repoName)
	if err != nil {
		return err
	}

	backupPath := m.checkForHomePath(absPath)
	id, ok := index.ContainPath(backupPath)
	if !ok {
		id, err = index.Add(backupPath)
		if err != nil {
			return err
		}
	}

	destPath := m.createDestPath(repoName, id)
	err = m.move(absPath, destPath)
	if err != nil {
		err = index.Remove(id)
		if err != nil {
			return err
		}
	}
	return err
}

func (m Mover) Remove(id string) (string, error) {
	for repo := range m.Config.Repository {
		index, err := m.checkIndex(repo)
		if err != nil {
			return "", err
		}
		if !index.CheckFile(id) {
			continue
		}
		err = m.RemoveFromRepository(repo, id)
		if err != nil {
			return "", err
		}
		return repo, nil
	}
	return "", errors.New("File not found.")
}

func (m Mover) RemoveFromRepository(repoName string, id string) error {
	index, err := m.checkIndex(repoName)
	if err != nil {
		return err
	}
	err = m.Git.Pull(repoName)
	if err != nil {
		return err
	}

	if !index.CheckFile(id) {
		return errors.New("File not found.")
	}
	err = m.Git.Remove(repoName, id)
	if err != nil {
		return err
	}
	err = index.Remove(id)
	return err
}

func (m Mover) List(repoName string) (map[string]string, error) {
	index, err := m.checkIndex(repoName)
	if err != nil {
		return nil, err
	}
	return index.Index, nil
}

func (m Mover) Update(repoName string) []error {
	var errs []error
	index, err := m.checkIndex(repoName)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	err = m.Git.Pull(repoName)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	for id, p := range index.Index {
		sourcePath := m.retrieveHomePath(p)
		destPath := m.createDestPath(repoName, id)
		err := m.Git.Remove(repoName, id)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		err = m.move(sourcePath, destPath)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (m Mover) Restore(repoName string) []error {
	index, err := m.checkIndex(repoName)
	var errs []error
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	err = m.Git.Pull(repoName)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	for id := range index.Index {
		err = m.RestoreFile(repoName, id)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (m Mover) RestoreFile(repoName string, id string) error {
	index, err := m.checkIndex(repoName)
	if err != nil {
		return err
	}
	restorePath := m.retrieveHomePath(index.Index[id])
	backupPath := m.createDestPath(repoName, id)
	err = m.move(backupPath, restorePath)
	return err
}

func (m Mover) move(source string, destination string) error {
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		err = copy.Copy(source, destination)
		if err != nil {
			return err
		}
	case mode.IsRegular():
		in, err := os.Open(source)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(destination)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		return err
	}

	return nil
}

func (m Mover) absolutePath(path string) (string, error) {
	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[1:]), nil
	}
	return filepath.Abs(path)
}

func (m Mover) checkForHomePath(path string) string {
	if strings.HasPrefix(path, "/home") {
		t := strings.Split(path, "/")
		t[2] = moverHome
		return strings.Join(t[2:], "/")
	}

	if strings.Contains(path, `\User`) {
		t := strings.Split(path, `\`)
		t[2] = moverHome
		return strings.Join(t[2:], `\`)
	}

	return path
}

func (m Mover) retrieveHomePath(path string) string {
	if strings.HasPrefix(path, "{HOME}") {
		home, _ := os.UserHomeDir()
		path = strings.Replace(path, moverHome, home, -1)
	}
	return path
}

func (m Mover) createDestPath(repoName string, id string) string {
	return path.Join(m.Path, repoName, id)
}
