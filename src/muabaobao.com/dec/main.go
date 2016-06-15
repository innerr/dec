package main

import (
	"errors"
	"os"
	"sort"
	"ember/cli"
)

func main() {
	hub := cli.NewRpcHub(os.Args[1:], NewServer, &Client{}, "/")
	hub.Run()
}

type Client struct {
	Send func(tags Tags) (err error) `args:"tags"`
	Get func(tags Tags) (err error) `args:"tags"`
	//Gets func(tags Tags) (events Events, err error) `args:"tags" return:"events"`
}

func (p *Server) Send(tags Tags) (err error) {
	sort.Strings(tags)
	p.cache[key(tags)] = tags
	return
}

func (p *Server) Get(tags Tags) (err error) {
	sort.Strings(tags)
	if p.cache[key(tags)] == nil {
		err = ErrEventNotExists
	}
	return
}

func NewServer(args []string) (p interface{}, err error) {
	p = &Server{make(map[string]Tags)}
	return
}

type Server struct {
	cache map[string]Tags
}

func key(tags Tags) (str string) {
	str = "."
	for _, tag := range tags {
		str += tag + "."
	}
	return
}

type Events []Tags

type Tags []string

var ErrEventNotExists = errors.New("event not exists")
