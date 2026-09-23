package hashtree

import "crypto/sha256"

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
