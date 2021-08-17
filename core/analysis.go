package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"regexp"

	"gopkg.in/yaml.v2"
)

func GetImages(inputPath string) []string {
	file, err := os.Stat(inputPath)
    if err != nil {
		fmt.Printf("Can not access < %s >. No such file or directory", inputPath)
        os.Exit(0)
    }

	result := []string{}
    if file.IsDir() {
		YAMLs, _ := GetAllYAML(inputPath)
		
		for _, YAML := range YAMLs {
			cont := ReadFile(YAML)
			if ! IsYAML(cont) {
				fmt.Printf("File < %s > is not a valid YAML file", YAML)
				continue
			}
			result = append(result, GetImagesFromContent(cont)...)
		}
	} else {
		cont := ReadFile(inputPath)
		if !IsYAML(cont) {
			fmt.Printf("File < %s > is not a valid YAML file", inputPath)
			os.Exit(0)
		}

		result = GetImagesFromContent(cont)
	}

	if len(result) == 0 {
		fmt.Println("No image found.")
		os.Exit(0)
	} else {
		fmt.Println("Images:")
		for _, image := range result {
			fmt.Println("  " + image)
		}
	}

	return result
}

func GetAllYAML(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() {
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllYAML(dirPth + PthSep + fi.Name())
		} else {
			ok := strings.HasSuffix(strings.ToLower(fi.Name()), ".yaml")
			if ok {
				files = append(files, dirPth + PthSep + fi.Name())
			}
		}
	}

	for _, table := range dirs {
		temp, _ := GetAllYAML(table)
		files = append(files, temp...)
	}

	return files, nil
}

func IsYAML(cont string) bool {
	return yaml.Unmarshal([]byte(cont), make(map[string]interface{})) == nil
}

func GetImagesFromContent(cont string) []string {
	lines := strings.Split(cont, "\n")
	pattern, _ := regexp.Compile("[\\s'\"]")

	result := make([]string, 0)
	for _, line := range lines {
		res := strings.Contains(line, "image:")
		if res {
			splited := strings.Split(line, "image:")
			result = append(result, pattern.ReplaceAllString(splited[len(splited)-1], ""))
		}
	}

	return result
}

func ReadFile(filePath string) string {
	if !IsFileExists(filePath) {
		fmt.Printf("File < %s > does not exist", filePath)
		os.Exit(0)
	}

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("read file fail", err)
		return ""
	}
	defer f.Close()

	fd, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("read to fd fail", err)
		return ""
	}

	return string(fd)
}