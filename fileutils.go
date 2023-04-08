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
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// while this is a definite progress bar,
// you can find an indefinite one here: https://github.com/alcionai/corso/blob/e09c12077847389b04745a8dd11c73b6162e2767/src/internal/observe/observe.go#L203-L211
// ^ saved in archive.is and archive.org/web
func newBar(maxBytes int64, name string) (*mpb.Progress, *mpb.Bar) {
	p := mpb.New(mpb.WithWidth(32))

	bar := p.New(maxBytes,
		mpb.BarStyle().Lbound("").Filler("█").Tip(" ").Padding(" ").Rbound(""),
		mpb.PrependDecorators(
			// len(name) + 1 ensure the name has one space on the right, no idea what decor.DidentRight does
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
			// upon completing, replace ETA decorator with "done" message
			// (actually not needed as we are removing bar on complete, but keeping it for reference)
			decor.OnComplete(
				decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 4}), "done",
			),
		),
		mpb.AppendDecorators(decor.CountersKibiByte("% .2f / % .2f")),
		mpb.BarRemoveOnComplete(),
	)

	return p, bar
}

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
func downloadFile(url, toLocation string) (int, error) {
	dirName := filepath.Dir(toLocation)
	fileName := filepath.Base(toLocation)
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		log.Println("Unable to create folder", dirName, "for file", fileName)
		return exitDirNotModifiable, err
	}

	// TODO check if it can successfully override existing file
	out, err := os.Create(toLocation)
	if err != nil {
		return exitFileNotModifiable, err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return exitUrlNotDownloadable, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return exitCannotDetermineUserHomeDir, fmt.Errorf("bad status: %s for %s", resp.Status, url)
	}

	p, bar := newBar(
		resp.ContentLength,
		"↓ "+fileName,
	)
	r := bar.ProxyReader(resp.Body)
	defer r.Close()

	_, err = io.Copy(out, r)
	p.Wait()
	if err != nil {
		return exitFileNotModifiable, err
	}

	return 0, nil
}

func markFileAsExecutableOnUnix(fullPath string) {
	cmd := exec.Command("chmod", "+x", fullPath)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("cannnot mark file as executable:", fullPath)
		panic(err)
	}
}
