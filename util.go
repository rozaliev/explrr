package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func SetValue(path string, attr string, value string) error {
	return ioutil.WriteFile(path+"/"+attr, []byte(value), 0644)
}

func GetValue(path string, attr string) (string, error) {
	data, err := ioutil.ReadFile(path + "/" + attr)
	if err != nil {
		return "", err
	}

	return strings.Trim(string(data), " \n\t"), nil
}

func FindDevice(path, attr, val string) (string, error) {
	dirs, err := listDir(path)
	if err != nil {
		return "", err
	}

	for _, dir := range dirs {
		data, err := GetValue(dir, attr)
		if err != nil {
			continue
		}

		if data == val {
			return dir, nil
		}
	}

	return "", errors.New("device was not found")
}

func listDir(root string) ([]string, error) {
	res := []string{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if root == path {
			return nil
		}
		res = append(res, path)

		if info.IsDir() {
			return filepath.SkipDir
		}

		return nil
	})
	return res, err
}
