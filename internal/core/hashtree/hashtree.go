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

type Hasher struct {
	hash func(left, right []byte) []byte
	size int
}

type node struct {
	hash   []byte
	parent *node
}

func SHA256() Hasher {
	return Hasher{
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

func FromBase(base []byte, h Hasher) (*HashTree, error) {
	if len(base) == 0 {
		return nil, errors.New("base is empty")
	}

	if len(base)%h.size != 0 {
		return nil, errors.New("invalid base length")
	}

	terminals := make([]*node, len(base)/h.size)
	for i := range terminals {
		terminals[i] = &node{hash: base[i*h.size : (i+1)*h.size]}
	}

	current := terminals
	for len(current) > 1 {
		next := make([]*node, (len(current)+1)/2)

		for i := range next {
			left := current[i*2]
			var right *node

			if i*2+1 < len(current) {
				right = current[i*2+1]
			} else {
				right = left
			}

			next[i] = &node{hash: h.hash(left.hash, right.hash)}
			left.parent = next[i]
			right.parent = next[i]
		}

		current = next
	}

	tree := &HashTree{
		root:      current[0],
		terminals: terminals,
		hasher:    h,
	}

	return tree, nil
}
