package main

import (
	"fmt"
	"ember/http/rpc"
	"ember/cli"
)

func (p *Cmds) Get(args []string) {
	ret, err := p.client.Call("Get", args)
	cli.Check(err)
	fmt.Println(ret[0])
}

func (p *Cmds) _Gets(args []string) {
	ret, err := p.client.Call("Gets", args)
	cli.Check(err)
	for _, it := range ret {
		fmt.Println(it)
	}
}

func (p *Cmds) Gets(args []string) {
	ret, err := p.client.Call("Gets", args)
	cli.Check(err)
	for _, it := range ret[0].([]Event) {
		fmt.Printf("%s %s %s\n", it.Path, it.Key, it.Event)
	}
}


func (p *Cmds) Set(args []string) {
	_, err := p.client.Call("Set", args)
	cli.Check(err)
}

type Cmds struct {
	client *rpc.Client
}
