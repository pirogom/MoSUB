package main

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

/**
*	decodeBase64
**/
func decodeBase64(b64 string) []byte {
	decodeStr, b64err := base64.StdEncoding.DecodeString(b64)
	if b64err == nil {
		return decodeStr
	}
	return nil
}

/**
*	splitFilename
**/
func splitFilename(fname string) (name string, ext string, err error) {

	lastDot := strings.LastIndex(fname, ".")

	if lastDot == -1 {
		return "", "", errors.New("bad file name")
	}

	name = fname[:lastDot]
	ext = fname[lastDot+1:]

	return name, ext, nil
}

/**
*	isWindows
**/
func isWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

/**
*	isDarwin
**/
func isDarwin() bool {
	if runtime.GOOS == "darwin" {
		return true
	}
	return false
}

/**
*	isExistFile
**/
func isExistFile(fname string) bool {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return false
	}
	return true
}

/**
*	fixWorkingDirectory
**/
func fixWorkingDirectory() {
	ex, err := os.Executable()
	if err != nil {
		os.Exit(0)
		return
	}
	exPath := filepath.Dir(ex)

	workDir := filepath.Join(exPath, "work")

	if !isExistFile(workDir) {
		os.Mkdir(workDir, 0644)
	}

	currDir, cdErr := os.Getwd()

	if cdErr != nil {
		os.Exit(0)
		return
	}

	if currDir != exPath {
		chErr := os.Chdir(exPath)

		if chErr != nil {
			os.Exit(0)
			return
		}
	}
}

/**
*	workFilepath
**/
func workFilepath(fname string) string {
	return filepath.Join("work", fname)
}
