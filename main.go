package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/antelman107/net-wait-go/wait"
)

func main() {
	var proto string
	flag.StringVar(&proto, "proto", "tcp", "tcp")
	var addrs string
	flag.StringVar(&addrs, "addrs", "", "address:port")
	var deadlineMS uint
	flag.UintVar(&deadlineMS, "deadline", 10000, "deadline in milliseconds")
	var delayMS uint
	flag.UintVar(&delayMS, "delay", 100, "delay in milliseconds")
	var debug bool
	flag.BoolVar(&debug, "debug", false, "debug messages toggler")

	flag.Parse()

	addrsSlice := strings.FieldsFunc(addrs, func(c rune) bool {
		return c == ','
	})

	if len(addrsSlice) == 0 {
		log.Println("addrs are not set")
		flag.Usage()

		os.Exit(1)
	}

	if wait.Do(
		proto,
		addrsSlice,
		time.Duration(delayMS)*time.Millisecond,
		time.Duration(deadlineMS)*time.Millisecond,
		debug,
	) {
		return
	}

	os.Exit(2)
}
