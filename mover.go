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

func CreateMover(config Config, git GitHelper, Path string) Mover {
	return Mover{
		Config: config,
		Git:    git,
		Path:   Path,
	}
}

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

	destPath := m.createDestPath(absPath, repoName)
	err = m.move(absPath, destPath)
	if err != nil {
		return err
	}

	// Handle Home directory
	// TODO: same for windows
	if strings.HasPrefix(absPath, "/home") {
		t := strings.Split(absPath, "/")
		t[2] = "{USER}"
		absPath = strings.Join(t, "/")
	}

	index, err := m.checkIndex(repoName)
	if err != nil {
		return err
	}

	if !index.ContainPath(absPath) {
		_, err = index.Add(absPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m Mover) Remove(repoName string, id string) error {
	index, err := m.checkIndex(repoName)
	if err != nil {
		return err
	}

	if index.CheckFile(id) {
		return errors.New("File not found.")
	}
	path := m.createDestPath(index.Index[id], repoName)
	err = os.RemoveAll(path)
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

	for _, v := range index.Index {
		if strings.Contains(v, "{USER}") {
			u, err := user.Current()
			if err != nil {
				return err
			}
			strings.Replace(v, "{USER}", u.Name, 1)
		}
		destPath := m.createDestPath(v, repoName)
		err = m.move(v, destPath)
		if err != nil {
			return err
		}
	}
	return nil
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

func (m Mover) createDestPath(repoName string, id string) string {
	return fmt.Sprint(m.Path, "/", repoName, "/", id)
}
