package clients

import (
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io/ioutil"
	"log"
	"os"
)

type GitParam struct {
	URL    string
	Folder string
}

func (param *GitParam) GetMemFS() (fs billy.Filesystem, err error) {
	fs = memfs.New()
	_, err = git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: param.URL,
	})
	return
}

func (param *GitParam) GetDirList(fs billy.Filesystem) (files []os.FileInfo, err error) {
	files, err = fs.ReadDir(param.Folder)
	return
}

func (param *GitParam) GetFileContent(fs billy.Filesystem, fileName string) (content string, err error) {
	file, err := fs.Open(param.Folder + string(os.PathSeparator) + fileName)
	if err != nil {
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	b, err := ioutil.ReadAll(file)
	content = string(b)
	return
}