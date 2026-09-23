package hashtree

import (
	"errors"
)

type node struct {
	hash   []byte
	parent *node
}

type HashTree struct {
	root      *node
	terminals []*node
	hasher    Hasher
}

func (tree *HashTree) RootHash() []byte {
	return tree.root.hash
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
