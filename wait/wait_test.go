package wait

import (
	"context"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func getTCPServer(proto, addr string, t *testing.T) io.Closer {
	srv, err := net.Listen(proto, addr)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			conn, err := srv.Accept()
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					t.Errorf("error listening UDP %s: %s", addr, err.Error())
				}
				return
			}

			_ = conn.Close()
		}
	}()

	return srv
}

func getUDPServer(proto, addr string, t *testing.T) io.Closer {
	udpAddr, err := net.ResolveUDPAddr(proto, addr)
	if err != nil {
		t.Errorf("error resolving UDP address %s: %s", udpAddr, err.Error())
		return nil
	}

	conn, err := net.ListenUDP(proto, udpAddr)
	if err != nil {
		t.Errorf("error listening UDP %s: %s", udpAddr, err.Error())
		return nil
	}

	go func() {
		for {
			var buf [1]byte
			_, udpRemoteAddr, err := conn.ReadFromUDP(buf[0:])
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					t.Errorf("error listening UDP %s: %s", addr, err.Error())
				}
				return
			}

			_, err = conn.WriteToUDP(buf[0:], udpRemoteAddr)
			if err != nil {
				t.Errorf("error writing from server to UDP client: %s", err.Error())
				return
			}
		}
	}()

	return conn
}

func TestDO(t *testing.T) {
	ctx := context.Background()

	type data struct {
		name          string
		addr          string
		reqAddr       string
		proto         string
		packet        string
		result        bool
		contextCancel bool
	}

	for _, row := range []data{
		{
			name:    "tcp success",
			addr:    "localhost:6432",
			reqAddr: "localhost:6432",
			proto:   "tcp",
			result:  true,
		},
		{
			name:    "tcp fail",
			addr:    "localhost:6432",
			reqAddr: "localhost:6431",
			proto:   "tcp",
			result:  false,
		},
		{
			name:    "udp success",
			addr:    "localhost:6433",
			reqAddr: "localhost:6433",
			proto:   "udp",
			packet:  "1",
			result:  true,
		},
		{
			name:    "udp fail",
			addr:    "localhost:6434",
			reqAddr: "localhost:6435",
			proto:   "udp",
			packet:  "1",
			result:  false,
		},
		{
			name:          "context cancel",
			addr:          "localhost:6433",
			reqAddr:       "localhost:6433",
			proto:         "udp",
			packet:        "1",
			result:        false,
			contextCancel: true,
		},
	} {
		r := row
		t.Run(row.name, func(t *testing.T) {
			var srv io.Closer
			if r.proto == "udp" {
				srv = getUDPServer(r.proto, r.addr, t)
			} else {
				srv = getTCPServer(r.proto, r.addr, t)
			}
			defer srv.Close()

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			if row.contextCancel {
				cancel()
			}

			e := New(
				WithProto(r.proto),
				WithUDPPacket([]byte(r.packet)),
				WithDebug(false),
				WithDeadline(time.Second*2),
				WithContext(ctx),
			)

			if e.Do([]string{r.reqAddr}) != r.result {
				t.Errorf("%s result is not %#v", r.name, r.result)
			}
		})
	}
}
