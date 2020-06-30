package wait

import (
	"net"
	"testing"
	"time"
)

func getServer(addr string) net.Listener {
	srv, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			conn, _ := srv.Accept()
			if conn != nil {
				_ = conn.Close()
			}
		}
	}()

	return srv
}

func TestTCP(t *testing.T) {
	ok := "localhost:6432"
	notok := "localhost:6431"

	srv := getServer(ok)
	defer srv.Close()

	if !Do("tcp", []string{ok}, time.Millisecond*100, time.Second, true) {
		t.FailNow()
	}

	if Do("tcp", []string{notok}, time.Millisecond*100, time.Second, true) {
		t.FailNow()
	}
}
