package main

func (p *Tree) SetValue(value string) {
	p.value = value
}

func (p *Tree) Value() string {
	return p.value
}

func (p *Tree) Children() map[string]*Tree {
	return p.children
}

func (p *Tree) Child(name string, create bool) *Tree {
	child, ok := p.children[name]
	if !ok && create {
		child = NewTree()
		p.children[name] = child
	}
	return child
}

func NewTree() *Tree {
	return &Tree{children: make(map[string]*Tree)}
}

type Tree struct {
	children map[string]*Tree
	value string
}
