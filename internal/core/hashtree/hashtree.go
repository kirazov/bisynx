package hashtree

type node struct {
	hash   []byte
	parent *node
}

type HashTree struct {
	top       *node
	terminals []*node
	hasher    Hasher
}

func (tree *HashTree) TopHash() []byte {
	return tree.top.hash
}

func FromHashes(hashes [][]byte, h Hasher) (*HashTree, error) {
	compiler := NewCompiler(h)
	for _, hash := range hashes {
		err := compiler.Add(hash)
		if err != nil {
			return nil, err
		}
	}
	return compiler.Compile()
}
