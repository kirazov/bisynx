package hashtree

import "errors"

type Compiler struct {
	hasher  Hasher
	current []*node
	count   int
	next    *Compiler
	pending *node
}

func NewCompiler(hasher Hasher) *Compiler {
	return &Compiler{hasher: hasher}
}

func join(left, right *node, h Hasher) *node {
	hash := h.Hash(left.hash, right.hash)
	parent := &node{hash: hash}
	left.parent = parent
	right.parent = parent
	return parent
}

func (c *Compiler) add(node *node) {
	c.count++
	if c.pending != nil {
		parent := join(c.pending, node, c.hasher)
		c.pending = nil

		if c.next == nil {
			c.next = NewCompiler(c.hasher)
		}
		c.next.add(parent)
	} else {
		c.pending = node
	}
}

func (c *Compiler) Add(hash []byte) {
	node := &node{hash: hash}
	c.current = append(c.current, node)
	c.add(node)
}

func (c *Compiler) compile() *node {
	if c.pending != nil && c.count > 1 {
		parent := join(c.pending, c.pending, c.hasher)
		c.pending = nil
		c.next.add(parent)
	}

	if c.next != nil {
		return c.next.compile()
	}

	return c.pending
}

func (c *Compiler) Compile() (*HashTree, error) {
	if c.count == 0 {
		return nil, errors.New("cannot compile empty tree")
	}

	root := c.compile()
	return &HashTree{
		root:      root,
		terminals: c.current,
		hasher:    c.hasher,
	}, nil
}
