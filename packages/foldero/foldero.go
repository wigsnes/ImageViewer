package foldero

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/wigsnes/imageViewer/packages/fileo"
)

type FolderInfo struct {
	Path          string
	Name          string
	NumberOfFiles int
	Column        string
}

func getFiles(path string) {

}

func GetFolderInfo(filePath, path string) []FolderInfo {
	var folders []FolderInfo

	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		folderPath := fmt.Sprintf("%s/%s", filePath, f.Name())
		folderFiles, _ := ioutil.ReadDir(folderPath)
		numFiles := fileo.NumberOfFiles(folderFiles)
		name := f.Name()
		folders = append(folders, FolderInfo{Path: path + f.Name(), Name: name, NumberOfFiles: numFiles, Column: "3"})
	}

	return folders
}
