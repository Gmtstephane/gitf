package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const ()

func main() {
	basePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}

	gitPath := basePath + "/git"

	if err := mkdirp(gitPath); err != nil {
		fmt.Println(err)
		return
	}

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		return
	}
	urlClone := argsWithoutProg[0]
	name := strings.TrimSuffix(urlClone, filepath.Ext(urlClone))
	parsedurl, err := url.Parse(name)
	if err != nil {
		return
	}

	//create hostpath folder if not exist
	hostpath := gitPath + "/" + parsedurl.Host
	splittedPath := strings.Split(parsedurl.Path, "/")
	if err := mkdirp(hostpath); err != nil {
		fmt.Println(err)
		return
	}
	currPath := hostpath
	for _, s := range splittedPath {
		currPath += "/" + s
		if err := mkdirp(currPath); err != nil {
			fmt.Println(err)
			return
		}
	}

	cmd := exec.Command("git", "clone", urlClone, currPath)

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	// Execute the command
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}

	// log.Println(stdBuffer.String())

}

func mkdirp(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
