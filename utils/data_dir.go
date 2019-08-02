package utils

import (
	"os"
	"path/filepath"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
)

// DataDir return a data directory for the current program
func DataDir(name string) (string, error) {
	if runtime.GOOS == "windows" {
		return dirWindows(name)
	} else if runtime.GOOS == "darwin" {
		return dirDarwin(name)
	}

	return dirUnix(name)
}

func dirWindows(name string) (string, error) {
	localAppData := os.Getenv("LOCALAPPDATA")

	if localAppData == "" {
		homedir, err := homedir.Dir()
		if err != nil {
			return "", err
		}

		localAppData = filepath.Join(homedir, "AppData", "Local")
	}

	dir := filepath.Join(localAppData, name)
	return dir, nil
}

func dirDarwin(name string) (string, error) {
	homedir, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(homedir, "Caches", name)
	return dir, nil
}

func dirUnix(name string) (string, error) {
	xdgCacheHome := os.Getenv("XDG_CACHE_HOME")

	if xdgCacheHome == "" {
		homedir, err := homedir.Dir()
		if err != nil {
			return "", err
		}

		xdgCacheHome = filepath.Join(homedir, ".cache")
	}

	dir := filepath.Join(xdgCacheHome, name)
	return dir, nil
}
