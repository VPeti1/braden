package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func InstallPackage(pkgName string) error {
	deps, err := GetDependencies(pkgName)
	if err != nil {
		return err
	}

	for _, dep := range deps {
		if !IsPackageInstalled(dep) {
			fmt.Printf("Installing dependency: %s\n", dep)
			err := InstallPackage(dep)
			if err != nil {
				return err
			}
		}
	}

	err = DownloadPackage(pkgName)
	if err != nil {
		return err
	}

	hash, err := GetPackageHash(pkgName)
	if err != nil {
		return err
	}

	if !VerifyFileHash(fmt.Sprintf("%s.pkg.tar.xz", pkgName), hash) {
		return fmt.Errorf("package verification failed for %s", pkgName)
	}

	err = ExtractPackage(pkgName)
	if err != nil {
		return err
	}

	version, err := GetPackageVersion(pkgName)
	if err != nil {
		return err
	}
	return SaveInstalledPackageVersion(pkgName, version)
}

func UpdatePackage(pkgName string) error {
	if !IsPackageInstalled(pkgName) {
		return InstallPackage(pkgName)
	}

	installedVersion, err := GetInstalledPackageVersion(pkgName)
	if err != nil {
		return err
	}

	remoteVersion, err := GetPackageVersion(pkgName)
	if err != nil {
		return err
	}

	if installedVersion != remoteVersion {
		fmt.Printf("Updating package %s from version %s to %s\n", pkgName, installedVersion, remoteVersion)
		err = InstallPackage(pkgName)
		if err != nil {
			return err
		}
	}

	return nil
}

func UninstallPackage(pkgName string) error {
	if !IsPackageInstalled(pkgName) {
		return fmt.Errorf("package %s is not installed", pkgName)
	}

	archivePath := fmt.Sprintf("%s.pkg.tar.xz", pkgName)
	err := os.Remove(archivePath)
	if err != nil {
		return err
	}

	installedFiles, err := ListInstalledFiles(pkgName)
	if err != nil {
		return err
	}

	for _, file := range installedFiles {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}

	versionFilePath := fmt.Sprintf("/usr/local/%s.version", pkgName)
	if err := os.Remove(versionFilePath); err != nil {
		return err
	}

	return nil
}

func ExtractPackage(pkgName string) error {
	archivePath := fmt.Sprintf("%s.pkg.tar.xz", pkgName)
	tarPath := fmt.Sprintf("%s.pkg.tar", pkgName)

	cmd := exec.Command("xz", "-d", archivePath)
	if err := cmd.Run(); err != nil {
		return err
	}

	tarFile, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	tarReader := tar.NewReader(tarFile)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		target := filepath.Join("/", header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}
		default:
			fmt.Printf("Unsupported file type: %c in %s\n", header.Typeflag, header.Name)
		}
	}

	if err := os.Remove(tarPath); err != nil {
		return err
	}

	return nil
}

func ListInstalledFiles(pkgName string) ([]string, error) {
	installedFiles := []string{}
	err := filepath.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(path, fmt.Sprintf("/usr/local/%s", pkgName)) {
			installedFiles = append(installedFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return installedFiles, nil
}

func IsPackageInstalled(pkgName string) bool {
	versionFilePath := fmt.Sprintf("/usr/local/%s.version", pkgName)
	if _, err := os.Stat(versionFilePath); os.IsNotExist(err) {
		return false
	}

	archivePath := fmt.Sprintf("%s.pkg.tar.xz", pkgName)
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		return false
	}

	return true
}
