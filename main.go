package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/antelman107/net-wait-go/wait"
)

func main() {
	conf, output, err := getFlags(os.Args[0], os.Args[1:])
	if err == flag.ErrHelp {
		fmt.Println(output)
		os.Exit(2)
	} else if err != nil {
		fmt.Println(err, "\n", output)
		os.Exit(1)
	}

	if wait.New(
		wait.WithProto(conf.proto),
		wait.WithWait(time.Duration(conf.delayMS)*time.Millisecond),
		wait.WithBreak(time.Duration(conf.breakMS)*time.Millisecond),
		wait.WithDeadline(time.Duration(conf.deadlineMS)*time.Millisecond),
		wait.WithDebug(conf.debug),
		wait.WithUDPPacket(conf.packetBytes),
	).Do(conf.addrsSlice) {
		return
	}

	os.Exit(1)
}
