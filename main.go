package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: braden <command> [args]")
		return
	}

	command := os.Args[1]
	initRepo()

	switch command {
	case "install":
		if len(os.Args) < 3 {
			fmt.Println("Usage: braden install <package>")
			return
		}
		packageName := os.Args[2]
		err := InstallPackage(packageName)
		if err != nil {
			fmt.Println("Error installing package:", err)
		} else {
			fmt.Println("Package installed successfully.")
		}
	case "update":
		if len(os.Args) < 3 {
			fmt.Println("Usage: braden update <package>")
			return
		}
		packageName := os.Args[2]
		err := UpdatePackage(packageName)
		if err != nil {
			fmt.Println("Error updating package:", err)
		} else {
			fmt.Println("Package updated successfully.")
		}
	case "uninstall":
		if len(os.Args) < 3 {
			fmt.Println("Usage: braden uninstall <package>")
			return
		}
		packageName := os.Args[2]
		err := UninstallPackage(packageName)
		if err != nil {
			fmt.Println("Error uninstalling package:", err)
		} else {
			fmt.Println("Package uninstalled successfully.")
		}
	default:
		fmt.Println("Unknown command:", command)
	}
}
