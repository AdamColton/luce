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
	k32 := Key32(k[0]) + Key32(k[1])<<8 + Key32(k[2])<<16 + Key32(k[3])<<24
	return k32
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
