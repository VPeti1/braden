package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var repoURL string

func readFirstLine(filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new scanner for the file
	scanner := bufio.NewScanner(file)

	// Read the first line
	if scanner.Scan() {
		repoURL = scanner.Text()
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func DownloadPackage(pkgName string) error {
	url := fmt.Sprintf("%s/%s.pkg.tar.xz", repoURL, pkgName)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download package: %s", resp.Status)
	}

	out, err := ioutil.TempFile("", "*.pkg.tar.xz")
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func GetDependencies(pkgName string) ([]string, error) {
	url := fmt.Sprintf("%s/%s.dep", repoURL, pkgName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get dependencies: %s", resp.Status)
	}

	var deps []string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		deps = append(deps, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return deps, nil
}

func GetPackageHash(pkgName string) (string, error) {
	url := fmt.Sprintf("%s/%s.sha256", repoURL, pkgName)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get package hash: %s", resp.Status)
	}

	hash, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func GetPackageVersion(pkgName string) (string, error) {
	url := fmt.Sprintf("%s/%s.version", repoURL, pkgName)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get package version: %s", resp.Status)
	}

	version, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(version), nil
}
