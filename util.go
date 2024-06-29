package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func CreateDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func VerifyFileHash(filename, expectedHash string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false
	}

	return fmt.Sprintf("%x", hash.Sum(nil)) == expectedHash
}

func SaveInstalledPackageVersion(pkgName, version string) error {
	file, err := os.Create(fmt.Sprintf("/usr/local/%s.version", pkgName))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(version)
	return err
}

func GetInstalledPackageVersion(pkgName string) (string, error) {
	file, err := os.Open(fmt.Sprintf("/usr/local/%s.version", pkgName))
	if err != nil {
		return "", err
	}
	defer file.Close()

	version, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(version), nil
}

func CompareVersions(v1, v2 string) int {
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	for i := 0; i < len(v1Parts) && i < len(v2Parts); i++ {
		if v1Parts[i] < v2Parts[i] {
			return -1
		}
		if v1Parts[i] > v2Parts[i] {
			return 1
		}
	}

	if len(v1Parts) < len(v2Parts) {
		return -1
	}
	if len(v1Parts) > len(v2Parts) {
		return 1
	}

	return 0
}

func initRepo() {
	filePath := "/etc/braden.conf"
	err := readFirstLine(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
