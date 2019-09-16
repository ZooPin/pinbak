package pinbak

import (
	"errors"
	"fmt"
	"github.com/otiai10/copy"
	"os"
	"path/filepath"
	"strings"
)

type Mover struct {
	Config Config
	Git    GitHelper
	Path   string
	Index  *IndexManager
}

func (m Mover) Add(path string, repoName string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// Handle Home directory
	// TODO: same for windows
	if strings.HasPrefix(absPath, "/home") {
		t := strings.Split(absPath, "/")
		t[1] = "{USER}"
		absPath = strings.Join(t, "/")
	}

	// TODO: Check repoName exist
	destPath := m.createDestPath(absPath, repoName)
	err = copy.Copy(absPath, destPath)
	if err != nil {
		return err
	}

	_, err = m.Index.Add(absPath, repoName)
	if err != nil {
		return err
	}

	return nil
}

func (m Mover) Remove(repoName string, id string) error {
	// TODO: Don't work with Index struct
	var index Index
	index, err := index.GetIndex(m.Path, repoName)
	if err != nil {
		return err
	}
	if _, ok := index.Index[id]; !ok {
		return errors.New("File not found.")
	}

	path := m.createDestPath(index.Index[id], repoName)
	err = os.RemoveAll(path)
	if err != nil {
		return err
	}
	// TODO: Need refactoring for index manager all merge with Index struct
	return nil
}

func (m Mover) List(repoName string) (map[string]string, error) {
	var index Index
	index, err := index.GetIndex(m.Path, repoName)
	if err != nil {
		return nil, err
	}
	return index.Index, nil
}

func (m Mover) Update() error {
	return nil
}

func (m Mover) createDestPath(path string, repoName string) string {
	return fmt.Sprint(path, "/", repoName, "/", filepath.Base(path))
}
