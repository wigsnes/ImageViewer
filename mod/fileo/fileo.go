package fileo

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"stringo"
)

func OpenFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func CreateFile(name string) *os.File {
	file, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func NextWord(file *os.File) func() string {
	scanner := bufio.NewScanner(file)

	return func() string {
		if scanner.Scan() {
			return scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		return ""
	}
}

func Sort(file *os.File) {

}

func DeleteDuplicates(path string) {
	file := OpenFile(path)
	defer file.Close()
	newFile := CreateFile(string([]rune(file.Name())[:len(file.Name())-4]) + "_deldup.txt")
	defer newFile.Close()

	next := NextWord(file)
	dup := ""

	for line := next(); line != ""; line = next() {
		if dup != line {
			newFile.WriteString(line + "\n")
			dup = line
		}
	}
}

func DeleteForeignCharacters(path string) {
	file := OpenFile(path)
	defer file.Close()
	newFile := CreateFile(string([]rune(file.Name())[:len(file.Name())-4]) + "_delfor.txt")
	defer newFile.Close()

	characters := []string{"æ", "Æ", "ø", "Ø", "å", "Å"}
	contains := false
	next := NextWord(file)
	for line := next(); line != ""; line = next() {
		for _, character := range characters {
			contains = strings.Contains(line, character)
			if contains {
				break
			}
		}
		if !contains {
			newFile.WriteString(line + "\n")
		}
		contains = false
	}
}

func RemoveFile(filePath string) error {
	split := strings.Split(filePath, "/")
	fileName := split[len(split)-1]
	path := strings.Replace(filePath, fileName, "", 1)
	trash := fmt.Sprintf("%sRECYCLE.BIN/%s", path, fileName)

	if _, err := os.Stat(fmt.Sprintf("%sRECYCLE.BIN", path)); os.IsNotExist(err) {
		os.Mkdir(fmt.Sprintf("%sRECYCLE.BIN", path), 0777)
	}

	return os.Rename(filePath, trash)
}

func NumberOfFiles(files []fs.DirEntry) int {
	numFiles := 0
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if stringo.StringInSlice(f.Name(), []string{".jpg", ".jpeg", ".png", ".gif", ".webm", ".mp4", ".mov"}) {
			numFiles++
		}
	}
	return numFiles
}

// ReplaceCharacter - replace all characters of type from in file and replace with to.
func ReplaceCharacterInFileName(filePath, from, to string) error {
	split := strings.Split(filePath, "/")
	fileName := split[len(split)-1]
	path := split[:len(split)-1]
	if strings.Contains(fileName, from) {
		newFileName := strings.Replace(fileName, from, to, -1)
		newFilePath := fmt.Sprintf("%s%s", path, newFileName)
		return os.Rename(filePath, newFilePath)
	}
	return nil
}
