package hashtree

import (
	"crypto/sha256"
	"errors"
)

type HashTree struct {
	root      *node
	terminals []*node
	hasher    Hasher
}

type Hasher interface {
	Hash(left, right []byte) []byte
	Size() int
}

type funcHasher struct {
	hash func(left, right []byte) []byte
	size int
}

func (h *funcHasher) Hash(left, right []byte) []byte {
	return h.hash(left, right)
}

func (h *funcHasher) Size() int {
	return h.size
}

func SHA256() Hasher {
	return &funcHasher{
		hash: func(left, right []byte) []byte {
			buf := make([]byte, 0, len(left)+len(right))
			buf = append(buf, left...)
			buf = append(buf, right...)
			hash := sha256.Sum256(buf)
			return hash[:]
		},
		size: sha256.Size,
	}
}

type node struct {
	hash   []byte
	parent *node
}

func getParent(left, right *node, h Hasher) *node {
	hash := h.Hash(left.hash, right.hash)
	parent := &node{hash: hash}
	left.parent = parent
	right.parent = parent
	return parent
}

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

func (c *Compiler) add(node *node) {
	c.count++
	if c.pending != nil {
		parent := getParent(c.pending, node, c.hasher)
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
		parent := getParent(c.pending, c.pending, c.hasher)
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

func FromBase(base []byte, h Hasher) (*HashTree, error) {
	if len(base) == 0 {
		return nil, errors.New("base is empty")
	}

	size := h.Size()

	if len(base)%size != 0 {
		return nil, errors.New("invalid base length")
	}

	compiler := NewCompiler(h)
	for i := 0; i < len(base); i += size {
		compiler.Add(base[i : i+size])
	}

	return compiler.Compile()
}
