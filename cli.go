package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"flag"
	"strings"
)

type flags struct {
	proto        string
	addrs        string
	addrsSlice   []string
	deadlineMS   uint
	delayMS      uint
	breakMS      uint
	debug        bool
	packetBase64 string
	packetBytes  []byte
}

var (
	ErrFlagsNotSet = errors.New("addrs are not set")
)

func getFlags(binaryName string, args []string) (*flags, string, error) {
	flagSet := flag.NewFlagSet(binaryName, flag.ContinueOnError)
	var buf bytes.Buffer
	flagSet.SetOutput(&buf)

	var conf = &flags{}
	flagSet.StringVar(&conf.proto, "proto", "tcp", "tcp")
	flagSet.StringVar(&conf.addrs, "addrs", "", "address:port(,address:port,address:port,...)")
	flagSet.UintVar(&conf.deadlineMS, "deadline", 10000, "deadline in milliseconds")
	flagSet.UintVar(&conf.delayMS, "wait", 100, "delay of single request in milliseconds")
	flagSet.UintVar(&conf.breakMS, "delay", 50, "break between requests in milliseconds")
	flagSet.BoolVar(&conf.debug, "debug", false, "debug messages toggler")
	flagSet.StringVar(&conf.packetBase64, "packet", "", "UDP packet to be sent")

	err := flagSet.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}

	conf.addrsSlice = strings.FieldsFunc(conf.addrs, func(c rune) bool {
		return c == ','
	})

	if len(conf.addrsSlice) == 0 {
		return nil, "", ErrFlagsNotSet
	}

	conf.packetBytes, err = base64.StdEncoding.DecodeString(conf.packetBase64)
	if err != nil {
		return nil, "", err
	}

	return conf, buf.String(), nil
}
