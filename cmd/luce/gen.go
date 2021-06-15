package main

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"

	"github.com/adamcolton/luce/tools/key"
	"github.com/urfave/cli"
)

func keyCmd(c *cli.Context) error {
	k := key.New(0)
	fmt.Println(k.Code())
	return nil
}

func randCmd(c *cli.Context) error {
	seed := make([]byte, 8)
	cryptorand.Read(seed)

	rand.Seed(
		(int64(seed[0]) << (8 * 0)) |
			(int64(seed[1]) << (8 * 1)) |
			(int64(seed[2]) << (8 * 2)) |
			(int64(seed[3]) << (8 * 3)) |
			(int64(seed[4]) << (8 * 4)) |
			(int64(seed[5]) << (8 * 5)) |
			(int64(seed[6]) << (8 * 6)) |
			(int64(seed[7]) << (8 * 7)),
	)

	max := c.Int64("n")
	if b := c.Int("b"); b > 0 {
		fmt.Println(b)
		max = 1 << uint(b)
	}
	fmt.Println(rand.Int63n(max))
	return nil
}

func randBase64(c *cli.Context) error {
	b := make([]byte, c.Int("b"))
	cryptorand.Read(b)

	fmt.Println(base64.URLEncoding.EncodeToString(b))
	return nil
}
