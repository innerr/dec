package main

import (
	"bufio"
	"errors"
	"io"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TODO: open file cache
// TODO: file writing transaction

func (p *Persist) Save(path string, key string, event string) (err error) {
	path = p.path + SEP + path
	err = os.MkdirAll(path, 0700)
	if err != nil && !os.IsExist(err) {
		return errors.New("create dir failed while saving: " + err.Error())
	}

	file, err := os.OpenFile(path + SEP + DATA_FILE, os.O_RDWR | os.O_APPEND | os.O_CREATE | os.O_SYNC, 0640)
	if err != nil {
		return errors.New("open failed while saving: " + err.Error())
	}
	_, err = fmt.Fprintf(file, "%s" + LINE_SEP + "%s\n", key, event)
	return
}

func (p *Persist) Load(base string, loaded map[string]bool) (err error) {
	base = p.path + SEP + base
	err = filepath.Walk(base, func(abs string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file := filepath.Base(abs)
		if file != DATA_FILE {
			return nil
		}
		path := filepath.Dir(abs)
		if len(path) <= len(p.path) {
			return errors.New("invalid base: " + path + " vs " + p.path)
		}
		rel := path[len(p.path) + 1: len(path)]
		if loaded[rel] {
			return nil
		}
		err = p.load(abs, rel)
		if err == nil {
			loaded[rel] = true
		}
		return err
	})

	if os.IsNotExist(err) {
		err = nil
	}
	if err != nil {
		err = errors.New("walking failed while loading: " + err.Error())
	}
	return
}

func (p *Persist) load(abs string, rel string) (err error) {
	paths := strings.Split(rel, SEP)
	node := p.root
	for i := 0; i < len(paths); i++ {
		node = node.Child(paths[i], true)
	}

	file, err := os.Open(abs)
	if err != nil {
		return errors.New("open file failed while loading: " + err.Error())
	}

	r := bufio.NewReader(file)
	for {
		data, oversize, err := r.ReadLine()
		line := string(data)
		if oversize {
			return errors.New("invalid line, too long: " + line)
		}
		if err != nil {
			if err != io.EOF {
				return errors.New("read line failed: " + err. Error())
			} else {
				return nil
			}
		}
		fields := strings.Split(line, LINE_SEP)
		if len(fields) != 2 {
			return errors.New("invalid line: " + line)
		}
		node.SetValue(fields[0], fields[1])
	}
	return
}

func NewPersist(root *Tree, path string) (p *Persist, err error) {
	p =&Persist{root, path}
	return
}

type Persist struct {
	root *Tree
	path string
}

var ErrNotPath = errors.New("not path")

const DATA_FILE = "events.dec"
const SEP = "/"
const LINE_SEP = "\t"
