package common

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// PathExists 判断路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ExtractStringReg(val string, pattern string) []map[string]string {

	re := regexp.MustCompile(pattern)
	match := re.FindAllStringSubmatch(val, -1)
	groupNames := re.SubexpNames()
	resultList := make([]map[string]string, 0)
	for i := 0; i < len(match); i++ {
		result := make(map[string]string)
		if len(match[i]) == len(groupNames) {
			for j, name := range groupNames {
				if j != 0 && name != "" {
					result[name] = match[i][j]
				}
			}
		}
		resultList = append(resultList, result)
	}

	return resultList
}

func MkDirs(perm os.FileMode, dirs ...string) {
	for _, dir := range dirs {
		exists, _ := PathExists(dir)
		if !exists {
			err := os.MkdirAll(dir, perm)
			if err != nil {
				return
			}
		}
	}
}

func GetMimeType(f *os.File) string {
	buffer := make([]byte, 512)
	_, _ = f.Read(buffer)
	contentType := http.DetectContentType(buffer)
	return contentType
}

func DeCompress(tarFile, dest string) ([]string, error) {

	var fileList []string

	srcFile, err := os.Open(tarFile)
	if err != nil {
		return nil, err
	}
	defer srcFile.Close()

	mimeType := GetMimeType(srcFile)
	switch mimeType {
	case "application/x-gzip":
		srcFile.Seek(0, os.SEEK_SET)
		gr, err := gzip.NewReader(srcFile)
		if err != nil {
			return nil, err
		}
		defer gr.Close()
		fileList, err = deTarCompress(gr, dest)
	case "application/octet-stream":
		srcFile.Seek(0, os.SEEK_SET)
		fileList, err = deTarCompress(srcFile, dest)
	}

	return fileList, nil
}

func deTarCompress(srcFile io.Reader, dest string) ([]string, error) {
	fileList := make([]string, 0)
	tr := tar.NewReader(srcFile)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		filename := dest + "/" + hdr.Name
		file, err := createFile(filename)
		if err != nil {
			continue
		}
		io.Copy(file, tr)
		fileList = append(fileList, hdr.Name)
		file.Close()
	}

	return fileList, nil
}

func createFile(name string) (*os.File, error) {
	path := string([]rune(name)[0:strings.LastIndex(name, "/")])
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
