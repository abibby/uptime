package main

import (
	"os"
	"path"
)

func OpenFile(parts ...string) (*os.File, error) {

	p := path.Join(parts...)
	dir := path.Dir(p)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	return os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

}
