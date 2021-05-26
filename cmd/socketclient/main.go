package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/ds/bus/procbus"
)

func main() {
	inBus := iobus.ReaderConfig{
		CloseOnEOF: true,
	}.New(os.Stdin)
	in := procbus.Delim(inBus.In, '\n')
	addr, err := getSock(in)
	if err != nil {
		fmt.Println(err)
		return
	}
	if addr == "" {
		return
	}
	fmt.Print("  Connecting to", addr, "\n\n")
	conn, err := net.Dial("unix", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	go iobus.Writer(conn, in, nil)

	cr := iobus.ReaderConfig{
		CloseOnEOF: true,
	}.New(conn)
	for m := range cr.In {
		fmt.Print(string(m))
	}
}

func getSock(in <-chan []byte) (string, error) {
	local, err := filepath.Glob("*.sock")
	if err != nil {
		return "", err
	}

	tmp, err := filepath.Glob("/tmp/*.sock")
	if err != nil {
		return "", err
	}

	all := append(local, tmp...)
	if len(all) == 0 {
		fmt.Println("No sockets found")
		return "", nil
	}

	fmt.Println("  Sockets:")
	for i, s := range all {
		fmt.Printf("    %d\t%s\n", i, s)
	}
	fmt.Print("(socket) ")
	b := <-in
	idx, err := strconv.Atoi(string(b))
	if err == nil && idx < len(all) {
		return all[idx], nil
	}
	return "", nil
}
