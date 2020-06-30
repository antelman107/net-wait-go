package wait

import (
	"net"
	"testing"
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

	e := New(WithDebug(true))
	if !e.Do([]string{ok}) {
		t.FailNow()
	}

	if e.Do([]string{notok}) {
		t.FailNow()
	}
}
