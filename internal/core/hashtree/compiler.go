package hashtree

type Compiler struct {
	hasher     Hasher
	nodeBuffer []*node
	compiled   bool
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

func (c *Compiler) Add(hash []byte) error {
	if c.compiled {
		return ErrAlreadyCompiled
	}

	if len(hash) != c.hasher.Size() {
		return &HashSizeError{expect: c.hasher.Size(), got: len(hash)}
	}

	c.nodeBuffer = append(c.nodeBuffer, &node{hash: hash})
	return nil
}

func (c *Compiler) Compile() (*HashTree, error) {
	if c.compiled {
		return nil, ErrAlreadyCompiled
	}

	if len(c.nodeBuffer) == 0 {
		return nil, ErrEmptyTree
	}

	tree := &HashTree{terminals: c.nodeBuffer, hasher: c.hasher}
	for len(c.nodeBuffer) > 1 {
		nextBuffer := make([]*node, (len(c.nodeBuffer)+1)/2)

		for i := range nextBuffer {
			var left, right *node
			left = c.nodeBuffer[i*2]

			if i*2+1 < len(c.nodeBuffer) {
				right = c.nodeBuffer[i*2+1]
			} else {
				right = left
			}

			nextBuffer[i] = join(left, right, c.hasher)
		}

		c.nodeBuffer = nextBuffer
	}
	c.compiled = true

	tree.top = c.nodeBuffer[0]
	return tree, nil
}

func (c *Compiler) Reset() {
	c.nodeBuffer = nil
	c.compiled = false
}
