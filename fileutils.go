package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func openTargetFile(filename string) (io.WriteCloser, error) {
	if filename == "" {
		return nopWriteCloser{os.Stdout}, nil
	} else {
		return os.Create(filename)
	}
}

func openSourceFile(filename string) (io.ReadCloser, error) {
	if (filename == "") || (filename == "-") {
		return ioutil.NopCloser(os.Stdin), nil
	} else {
		return os.Open(filename)
	}
}

func chooseRendererBaseOnOutfile(filename string) (Renderer, error) {
	ext := filepath.Ext(filename)
	if ext == ".png" {
		return PngRenderer, nil
	} else if ext == ".svg" {
		return SvgRenderer, nil
	}

	return nil, errors.New("Unsupported extension: " + filename)
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }
