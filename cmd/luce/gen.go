package main

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/adamcolton/luce/tools/key"

	//"github.com/adamcolton/luce/util/"lrand
	"github.com/urfave/cli"
)

func keyCmd(c *cli.Context) error {
	k := key.New(0)
	fmt.Println(k.Code())
	return nil
}

func randCmd(c *cli.Context) error {
	max := c.Int64("n")
	if b := c.Int("b"); b > 0 {
		fmt.Println(b)
		max = 1 << uint(b)
	}
	rand.Seed(time.Now().UnixMicro())
	fmt.Println(rand.Int63n(max))
	return nil
}

func randBase64(c *cli.Context) error {
	b := make([]byte, c.Int("b"))
	cryptorand.Read(b)

	fmt.Println(base64.URLEncoding.EncodeToString(b))
	return nil
}

func rand32(c *cli.Context) error {
	rand.Seed(time.Now().UnixMicro())
	fmt.Println(rand.Uint32())
	return nil
}
