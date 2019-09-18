package pinbak

import (
	"errors"
	"fmt"
	"github.com/otiai10/copy"
	"io"
	"os"
	"os/user"
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

const moverUser = "{USER}"

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
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	index, err := m.checkIndex(repoName)
	if err != nil {
		return err
	}

	backupPath := absPath
	// Handle Home directory
	// TODO: same for windows & better way to handle this
	if strings.HasPrefix(absPath, "/home") {
		t := strings.Split(absPath, "/")
		t[2] = moverUser
		backupPath = strings.Join(t, "/")
	}
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

func (m Mover) Remove(repoName string, id string) error {
	index, err := m.checkIndex(repoName)
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

func (m Mover) Update(repoName string) error {
	index, err := m.checkIndex(repoName)
	if err != nil {
		return err
	}

	for id, path := range index.Index {
		sourcePath := m.handleHomePath(path)
		destPath := m.createDestPath(repoName, id)
		err := m.Git.Remove(repoName, id)
		if err != nil {
			return err
		}
		err = m.move(sourcePath, destPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m Mover) Restore(repoName string) error {
	index, err := m.checkIndex(repoName)
	if err != nil {
		return err
	}

	for id, _ := range index.Index {
		err = m.RestoreFile(repoName, id)
	}
	return err
}

func (m Mover) RestoreFile(repoName string, id string) error {
	index, err := m.checkIndex(repoName)
	if err != nil {
		return err
	}
	restorePath := m.handleHomePath(index.Index[id])
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

func (m Mover) handleHomePath(path string) string {
	if strings.HasPrefix(path, "/home") {
		u, _ := user.Current()
		path = strings.Replace(path, moverUser, u.Name, -1)
	}
	return path
}

func (m Mover) createDestPath(repoName string, id string) string {
	return fmt.Sprint(m.Path, "/", repoName, "/", id)
}
