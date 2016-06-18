package main

import (
	"errors"
	"os"
	"strings"
	"sync"
	"ember/cli"
)

func main() {
	hub := cli.NewRpcHub(os.Args[1:], NewServerByArgs, &Client{}, "/")
	hub.Run()
}

type Client struct {
	Set func(path string, key string, value string) (err error) `args:"path,key,event"`
	Get func(path string, key string) (event string, err error) `args:"path,key" return:"event"`
	Gets func(path string) (events []Event, err error) `args:"tags" return:"events"`
}

func (p *Server) Set(path string, key string, event string) (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	node, err := p.get(path, true)
	if err != nil {
		return
	}
	node.SetValue(key, event)
	err = p.persist.Save(path, key, event)
	return
}

func (p *Server) Get(path string, key string) (event string, err error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	node, err := p.get(path, false)
	if err != nil {
		return
	}
	event = node.GetValue(key)
	if len(event) == 0 {
		err = ErrEventNotExists
	}
	return
}

func (p *Server) Gets(path string) (events []Event, err error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	node, err := p.get(path, false)
	if err != nil {
		return
	}
	events = make([]Event, 0)
	events, err = p.gets(node, path, events)
	return
}

func (p *Server) gets(node *Tree, path string, events []Event) (result []Event, err error) {
	result = events
	for k, v := range node.Values() {
		result = append(result, Event{path, k, v})
	}
	for k, v := range node.Children() {
		result, err = p.gets(v, path + SEP + k, result)
	}
	return
}

func (p *Server) get(path string, writing bool) (node *Tree, err error) {
	err = p.check(path)
	if err != nil {
		return
	}

	if !writing {
		err = p.persist.Load(path, p.cached)
		if err != nil {
			return
		}
	}

	paths := strings.Split(path, SEP)
	node = p.cache
	for _, it := range paths {
		node = node.Child(it, writing)
		if node == nil {
			err = errors.New("path not exists: " + path)
			return
		}
	}
	return
}

func (p *Server) check(path string) (err error) {
	if strings.HasPrefix(path, SEP) || strings.HasSuffix(path, SEP) {
		err = errors.New("invalid path: " + path)
	}
	return
}

func NewServerByArgs(args []string) (p interface{}, err error) {
	if len(args) != 1 {
		err = errors.New("command args not matched")
		return
	}
	p, err = NewServer(args[0])
	return
}

func NewServer(path string) (p *Server, err error) {
	root := NewTree()
	persist, err := NewPersist(root, path)
	if err != nil {
		return
	}
	p = &Server{cache: root, persist: persist, cached: make(map[string]bool)}
	return
}

type Server struct {
	cache *Tree
	persist *Persist
	cached map[string]bool
	lock sync.Mutex
}

type Event struct {
	Path string
	Key string
	Event string
}

var ErrEventNotExists = errors.New("event not exists")
