package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"p2p/enc"
	"time"
)

type KeyFile struct {
	K string
	T string
}

func main() {
	var (
		argInterval string
		argFile     string
	)

	flag.StringVar(&argInterval, "i", "6h", "Time interval for which key will be valid. Accepts the following units: ns, us, ms, s, m, h\n")
	flag.StringVar(&argFile, "f", "keyfile.yaml", "Filename")
	flag.Parse()
	d, err := time.ParseDuration(argInterval)
	if err != nil {
		fmt.Printf("Failed to parse provided interval: %s\n", argInterval)
		os.Exit(1)
	}
	t := time.Now()
	t = t.Add(d)
	fmt.Printf("This key will be valid until %s\n", t.String())

	key := enc.MakeEncKey()

	file, err := os.Create(argFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var buf []byte
	buf = append(buf, key...)
	tb, _ := t.MarshalBinary()
	buf = append(buf, tb...)

	_, err = io.WriteString(file, string(buf))

	/*
		_, err = io.WriteString(file, string(key))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = io.WriteString(file, []byte(t.String()))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}*/

	file.Close()

	nf, err := os.OpenFile(argFile, os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("ERR %v", err)
	}
	var ret []byte
	_, err = nf.Read(ret)
	if err != nil {
		fmt.Printf("ERR %v", err)
	}
	fmt.Printf("LEN: %d", len(ret))
	newKey := ret[:32]
	if bytes.Compare(key, newKey) == 0 {
		fmt.Println("!!!")
	}
}
