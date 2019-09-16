package pinbak

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"os"
)

type IndexManager struct {
	Repository map[string]map[string]string `json:"repository"`
	Path       string                       `json:"-"`
}

type Index struct {
	Index map[string]string `json:"index"`
}

func (I Index) checkIndex(basePath string, repoName string) bool {
	path := fmt.Sprint(basePath, "/", repoName)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (I Index) GetIndex(basePath string, repoName string) (Index, error) {
	var index Index
	if !I.checkIndex(basePath, repoName) {
		index.Index = make(map[string]string)
	} else {
		path := fmt.Sprint(basePath, "/", repoName, "/.index")
		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			return Index{}, err
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return Index{}, err
		}
		err = json.Unmarshal(data, &index)
		if err != nil {
			return Index{}, err
		}
	}
	return index, nil
}

func (im *IndexManager) Add(path string, repoName string) (string, error) {
	var index Index
	index, err := index.GetIndex(im.Path, repoName)
	id := fmt.Sprintf("%s", uuid.NewV4())
	if !im.CheckRepository(repoName) {
		im.Repository[repoName] = make(map[string]string)
	}
	im.Repository[repoName][id] = path

	err = im.save(repoName)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (im *IndexManager) Remove(repoName string, id string) error {
	if !im.CheckRepository(repoName) {
		return errors.New("Repository name not found")
	}

	if !im.CheckFile(repoName, id) {
		return errors.New("File not found")
	}

	delete(im.Repository[repoName], id)

	err := im.save(repoName)
	if err != nil {
		return err
	}

	return nil
}

func (im IndexManager) CheckRepository(repoName string) bool {
	_, ok := im.Repository[repoName]
	return ok
}

func (im IndexManager) CheckFile(repoName string, id string) bool {
	_, ok := im.Repository[repoName][id]
	return ok
}

func (im IndexManager) save(repoName string) error {
	file, err := json.MarshalIndent(im, "", "  ")
	if err != nil {
		return err
	}
	path := fmt.Sprint(im.Path, repoName, "/.index")
	err = ioutil.WriteFile(path, file, 0644)
	if err != nil {
		return err
	}
	return nil
}
