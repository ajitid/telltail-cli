package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func fileOrFolderExists(fullpath string) bool {
	_, err := os.Stat(fullpath)
	return !errors.Is(err, fs.ErrNotExist)
}

func removeFileIfPresent(fullpath string) error {
	if fileOrFolderExists(fullpath) {
		if err := os.Remove(fullpath); err != nil {
			return cli.Exit("Unable to remove "+filepath.Base(fullpath)+" from path "+filepath.Dir(fullpath), exitFileNotModifiable)
		}
	}
	return nil
}

// TODO rename folder to dir
func removeFolderIfPresent(fullpath string) error {
	if fileOrFolderExists(fullpath) {
		if err := os.RemoveAll(fullpath); err != nil {
			return cli.Exit("Unable to remove "+fullpath, exitDirNotModifiable)
		}
	}
	return nil
}

func isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // It's better to use f.Readdirnames(n) because its implementation just reads names, while f.Readdir(n) reads each found file statistics as well.
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func removeDirIfEmpty(fullpath string) error {
	if empty, _ := isDirEmpty(fullpath); empty {
		return removeFolderIfPresent(fullpath)
	}

	return nil
}

func cmdExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// taken from https://stackoverflow.com/questions/11692860/how-can-i-efficiently-download-a-large-file-using-go
func downloadFile(url, toLocation string) (error, int) {
	dirName := filepath.Dir(toLocation)
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		fileName := filepath.Base(toLocation)
		log.Println("Unable to create folder", dirName, "for file", fileName)
		return err, exitDirNotModifiable
	}

	// TODO check if it can successfully override existing file
	out, err := os.Create(toLocation)
	if err != nil {
		return err, exitFileNotModifiable
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err, exitUrlNotDownloadable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s for %s", resp.Status, url), exitUrlNotDownloadable
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err, exitFileNotModifiable
	}

	return nil, 0
}

func markFileAsExecutableOnUnix(fullPath string) {
	cmd := exec.Command("chmod", "+x", fullPath)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("cannnot mark file as executable:", fullPath)
		panic(err)
	}
}
