package main

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	array := make([]byte, 0)
	err := walk(out, path+"\\", "", printFiles)
	if err != nil {
		return err
	}
	out.Write(array)
	return nil
}

func walk(out io.Writer, path string, prefix string, printFiles bool) error {
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	if !printFiles {
		fileInfo = deleteFilesFromSlice(fileInfo)
	}
	for index, file := range fileInfo {
		if index == len(fileInfo)-1 {
			addDataToWriter(out, file, prefix+"└───")
			if file.IsDir() {
				walk(out, path+file.Name()+"\\", prefix+"\t", printFiles)
			}
		} else {
			addDataToWriter(out, file, prefix+"├───")
			if file.IsDir() {
				walk(out, path+file.Name()+"\\", prefix+"│\t", printFiles)
			}
		}
	}
	return nil
}

func addDataToWriter(out io.Writer, fileInfo os.FileInfo, prefix string) {
	out.Write([]byte(prefix + getFormatedFileName(fileInfo)))
}

func getFormatedFileName(fileInfo os.FileInfo) string {
	fileSize := getFileSize(fileInfo)
	formatedName := fileInfo.Name()
	if fileSize == "" {
		return formatedName + "\n"
	}
	return formatedName + " " + getFileSize(fileInfo) + "\n"
}

func getFileSize(fileInfo os.FileInfo) string {
	if fileInfo.IsDir() {
		return ""
	}
	size := fileInfo.Size()
	if size > 0 {
		return formatFileSize(size)
	}
	return "(empty)"
}

func formatFileSize(size int64) string {
	return "(" + strconv.FormatInt(size, 10) + "b" + ")"
}

func deleteFilesFromSlice(arr []os.FileInfo) []os.FileInfo {
	resultSlice := make([]os.FileInfo, 0)
	for x, file := range arr {
		if file.IsDir() {
			resultSlice = append(resultSlice, arr[x])
		}
	}
	return resultSlice
}
