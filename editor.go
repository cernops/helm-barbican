package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const (
	DefaultEditor = "vim"
)

type Editor struct {
	Binary string
}

func NewEditor() (Editor, error) {
	bin := os.Getenv("EDITOR")
	if bin == "" {
		bin = DefaultEditor
	} else {
		elems := strings.Fields(bin)
		if len(elems) > 0 {
			bin = elems[0]
		}
	}

	path, err := exec.LookPath(bin)
	if err != nil {
		return Editor{Binary: bin}, err
	}
	return Editor{Binary: path}, nil
}

func (e Editor) Launch(path string) error {
	cmd := exec.Command(e.Binary, path)
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

func (e Editor) LaunchTemp(r io.Reader) ([]byte, string, error) {
	f, err := ioutil.TempFile("", "*.helm")
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

func randomString(n int) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return string(b)
}
