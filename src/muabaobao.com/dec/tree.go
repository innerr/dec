package main

func (p *Tree) Values() map[string]string {
	return p.values
}

func (p *Tree) SetValue(key string, value string) {
	p.values[key] = value
}

func (p *Tree) GetValue(key string) string {
	return p.values[key]
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
	return &Tree{make(map[string]*Tree), make(map[string]string)}
}

type Tree struct {
	children map[string]*Tree
	values map[string]string
}
