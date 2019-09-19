package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/xid"
	"io/ioutil"
	"os"
)

type Index struct {
	Index    map[string]string `json:"index"`
	Path     string            `json:"-"`
	RepoName string            `json:"-"`
	guid     xid.ID            `json:"-"`
}

func openIndex(basePath string, repoName string) (Index, error) {
	var index Index
	err := index.open(basePath, repoName)
	index.RepoName = repoName
	index.guid = xid.New()
	return index, err
}

func (I Index) checkIndex(basePath string, repoName string) bool {
	info, err := os.Stat(I.Path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (I *Index) open(basePath string, repoName string) error {
	I.Path = fmt.Sprint(basePath, "/", repoName, "/index")
	if !I.checkIndex(basePath, repoName) {
		I.Index = make(map[string]string)
		err := I.save()
		return err
	}

	f, err := os.Open(I.Path)
	defer f.Close()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &I)
	if err != nil {
		return err
	}
	return nil
}

func (I *Index) Add(path string) (string, error) {
	if I.Index == nil {
		I.Index = make(map[string]string)
	}
	id := fmt.Sprintf("%s", I.guid.String())
	I.Index[id] = path

	err := I.save()
	if err != nil {
		return "", err
	}

	return id, nil
}

func (I Index) CheckFile(id string) bool {
	_, ok := I.Index[id]
	return ok
}

func (I Index) ContainPath(path string) (string, bool) {
	for k, v := range I.Index {
		if v == path {
			return k, true
		}
	}
	return "", false
}

func (I Index) save() error {
	file, err := json.MarshalIndent(I, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(I.Path, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (I *Index) Remove(id string) error {
	if !I.CheckFile(id) {
		return errors.New("File not found")
	}
	delete(I.Index, id)

	err := I.save()
	if err != nil {
		return err
	}

	return nil
}
