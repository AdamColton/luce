package lhash

import (
	"hash"

	"github.com/adamcolton/luce/util/reflector"
)

// Hasher represents an object that can return a hash value representing itself.
type Hasher interface {
	Hash() []byte
}

type Hasher64 interface {
	Hash() uint64
}

type Hasher32 interface {
	Hash() uint32
}

type Factory func() hash.Hash

func TypeHasher[T any](f Factory) func(t T) []byte {
	return func(t T) []byte {
		h := f()
		h.Write(reflector.UnsafeByteSlice(t))
		return h.Sum(nil)
	}
}

type Factory64 func() hash.Hash64

func TypeHasher64[T any](f Factory64) func(t T) uint64 {
	return func(t T) uint64 {
		h := f()
		h.Write(reflector.UnsafeByteSlice(t))
		return h.Sum64()
	}
}

type Factory32 func() hash.Hash32

func TypeHasher32[T any](f Factory32) func(t T) uint32 {
	return func(t T) uint32 {
		h := f()
		h.Write(reflector.UnsafeByteSlice(t))
		return h.Sum32()
	}
}
