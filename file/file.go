package file

import (
	"log"
	"os"
	"path/filepath"
)

// Exists 检查目录或文件是否存在
func Exists(filepath string) bool {
	return checkExist(filepath)
}

func checkExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// List 获取指定目录含子目录下的所有文件, 返回列表
func List(path string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	var fileList []string
	for _, file := range files {
		if file.IsDir() {
			fileList = append(fileList, List(filepath.Join(path, file.Name()))...)
		} else {
			fileList = append(fileList, filepath.Join(path, file.Name()))
		}
	}
	return fileList
}
