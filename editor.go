package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	DefaultEditor = "vim"
)

type Editor struct {
	Binary string
}

func NewEditor() Editor {
	bin := os.Getenv("EDITOR")
	if bin == "" {
		bin = DefaultEditor
	}
	return Editor{Binary: bin}
}

func (e Editor) Launch(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	cmd := exec.Command(e.Binary, abs)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		if err, ok := err.(*exec.Error); ok {
			if err.Err == exec.ErrNotFound {
				return fmt.Errorf("unable to launch editor : %v", err)
			}
		}
		return fmt.Errorf("editor failed : %v", err)
	}
	return nil
}

func (e Editor) LaunchTemp(prefix string, suffix string, r io.Reader) ([]byte, string, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, "", err
	}
	defer os.Remove(f.Name())
	if _, err := io.Copy(f, r); err != nil {
		return nil, "", err
	}
	if err := e.Launch(f.Name()); err != nil {
		return nil, f.Name(), err
	}
	bytes, err := ioutil.ReadFile(f.Name())
	return bytes, f.Name(), err
}
