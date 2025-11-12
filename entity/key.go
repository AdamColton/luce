package entity

import (
	"crypto/rand"
	"hash/crc64"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/rye"
)

type Key []byte

var RandKeyLength = 8

func Rand() Key {
	k := make([]byte, RandKeyLength)
	rand.Read(k)
	return k
}

func Rand32() Key32 {
	k := make([]byte, 4)
	rand.Read(k)
	k32 := Key32(k[0]) + Key32(k[1])<<8 + Key32(k[2])<<16 + Key32(k[3])<<24
	return k32
}

func Rand64() Key64 {
	k := make([]byte, 8)
	rand.Read(k)
	return NewKey64(k)
}

var tab64 = crc64.MakeTable(crc64.ISO)

func (key Key) EntKey() Key {
	return key
}

func (key Key) Key32() Key32 {
	return Key32(rye.Deserialize.Uint32(key))
}

func (key Key) Hash64() uint64 {
	h64 := crc64.New(tab64)
	lerr.Must(h64.Write(key))
	return h64.Sum64()
}

type Key32 uint32

func (k32 Key32) EntKey() Key {
	var k [4]byte
	rye.Serialize.Uint32(k[:], uint32(k32))
	return k[:]
}

type Key64 uint64

func (k64 Key64) EntKey() Key {
	var k [4]byte
	rye.Serialize.Uint64(k[:], uint64(k64))
	return k[:]
}

func NewKey64(k Key) Key64 {
	k64 := Key64(k[0]) + Key64(k[1])<<8 + Key64(k[2])<<16 + Key64(k[3])<<24 + Key64(k[4])<<32 + Key64(k[5])<<40 + Key64(k[6])<<48 + Key64(k[7])<<56
	return k64
}
