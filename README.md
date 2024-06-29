# Braden Package Manager

A basic package manager written in Go. This package manager can handle installing, updating, and uninstalling packages.

## Features

- Install packages and their dependencies
- Update installed packages
- Uninstall packages
- Verify package integrity using SHA256 hashes

## Prerequisites

- Go (1.15 or later) (For compiling)
- `xz` command-line utility (for decompressing `.xz` archives)

### Install a package

braden install <package-name>

Example:

braden install neofetch

### Update a package

braden update <package-name>

Example:

braden update neofetch

Uninstall a package

### braden uninstall <package-name>

Example:

braden uninstall neofetch

### Repo configuration
    Create a file at /etc/braden.conf
    Put your repos url there 

### Example braden repo structure
    repo/ 
        neofetch-7.1.0.pkg.tar.xz
        neofetch-7.1.0.dep
        neofetch-7.1.0.sha256
        neofetch-7.1.0.version

### Warning 
This is just a proof of concept package manager 



