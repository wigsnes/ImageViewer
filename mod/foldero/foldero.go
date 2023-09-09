package foldero

import (
	"fmt"
	"log"
	"os"

	"fileo"
)

type FolderInfo struct {
	Path          string
	Name          string
	NumberOfFiles int
	Column        string
}

func getFiles(path string) {

}

func GetFolderInfo(filePath string) []FolderInfo {
	var folders []FolderInfo

	files, err := os.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		folderPath := fmt.Sprintf("%s/%s", filePath, f.Name())
		folderFiles, _ := os.ReadDir(folderPath)
		numFiles := fileo.NumberOfFiles(folderFiles)
		name := f.Name()
		folders = append(folders, FolderInfo{Path: filePath + f.Name(), Name: name, NumberOfFiles: numFiles, Column: "3"})
	}

	return folders
}
