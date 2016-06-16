package main

import (
	"errors"
	"os"
	"strings"
	"ember/cli"
)

func main() {
	hub := cli.NewRpcHub(os.Args[1:], NewServer, &Client{}, "/")
	hub.Run()
}

type Client struct {
	Send func(path string, event string) (err error) `args:"path,event"`
	Get func(path string) (event string, err error) `args:"path" return:"event"`
	Gets func(path string) (events []Event, err error) `args:"tags" return:"events"`
}

func (p *Server) Send(path string, event string) (err error) {
	node, err := p.get(path, true)
	if err != nil {
		return
	}
	node.SetValue(event)
	// TODO: save
	return
}

func (p *Server) Get(path string) (event string, err error) {
	node, err := p.get(path, false)
	if err != nil {
		return
	}
	event = node.Value()
	if event == "" {
		err = ErrEventNotExists
	}
	return
}

func (p *Server) Gets(path string) (events []Event, err error) {
	node, err := p.get(path, false)
	if err != nil {
		return
	}
	events = make([]Event, 0)
	events, err = p.gets(node, path, events)
	return
}

func (p *Server) gets(node *Tree, prefix string, events []Event) (result []Event, err error) {
	result = events
	for k, v := range node.Children() {
		if len(v.Children()) != 0 && v.Value() != "" {
			err = ErrNotPath
			return
		} else if v.Value() != "" {
			result = append(events, Event{prefix + "/" + k, v.Value()})
		} else {
			result, err = p.gets(v, prefix + "/" + k, result)
		}
	}
	return
}

func (p *Server) get(path string, create bool) (node *Tree, err error) {
	// TODO: load
	paths := strings.Split(path, "/")
	node = p.cache
	for i, it := range paths {
		node = node.Child(it, create)
		if node == nil {
			err = ErrPathNotExists
			return
		}
		if node.Value() != "" && i != len(paths) - 1 {
			err = ErrNotPath
			return
		}
	}
	return
}

func NewServer(args []string) (p interface{}, err error) {
	p = &Server{NewTree()}
	return
}

type Server struct {
	cache *Tree
}

type Event struct {
	Path string
	Event string
}

var ErrNotPath = errors.New("not path")
var ErrEventNotExists = errors.New("event not exists")
var ErrPathNotExists = errors.New("path not exists")
