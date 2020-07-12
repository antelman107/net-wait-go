package main

import (
	"encoding/base64"
	"flag"
	"fmt"
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
	flag.UintVar(&delayMS, "wait", 100, "delay of single request in milliseconds")

	var breakMS uint
	flag.UintVar(&breakMS, "delay", 50, "break between requests in milliseconds")

	var debug bool
	flag.BoolVar(&debug, "debug", false, "debug messages toggler")

	var packetBase64 string
	flag.StringVar(&packetBase64, "packet", "", "UDP packet to be sent")

	flag.Parse()

	addrsSlice := strings.FieldsFunc(addrs, func(c rune) bool {
		return c == ','
	})

	if len(addrsSlice) == 0 {
		log.Println("addrs are not set")
		flag.Usage()

		os.Exit(2)
	}

	packetBytes, err := base64.StdEncoding.DecodeString(packetBase64)
	if err != nil {
		fmt.Println("packet base64 decode error:", err)
		flag.Usage()

		os.Exit(2)
	}

	if wait.New(
		wait.WithProto(proto),
		wait.WithWait(time.Duration(delayMS)*time.Millisecond),
		wait.WithBreak(time.Duration(breakMS)*time.Millisecond),
		wait.WithDeadline(time.Duration(deadlineMS)*time.Millisecond),
		wait.WithDebug(debug),
		wait.WithUDPPacket(packetBytes),
	).Do(addrsSlice) {
		return
	}

	os.Exit(1)
}
