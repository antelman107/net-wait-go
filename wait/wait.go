package wait

import (
	"context"
	"log"
	"net"
	"sync"
	"time"
)

type Executor struct {
	Proto     string
	Wait      time.Duration
	Break     time.Duration
	Deadline  time.Duration
	Debug     bool
	UDPPacket []byte
	Context   context.Context
}

type Option func(*Executor)

func New(opts ...Option) *Executor {
	const (
		defaultProto     = "tcp"
		defaultWait      = 200 * time.Millisecond
		defaultBreak     = 50 * time.Millisecond
		defaultDeadline  = 15 * time.Second
		defaultDebug     = false
		defaultUDPPacket = ""
	)

	e := &Executor{
		Proto:     defaultProto,
		Wait:      defaultWait,
		Break:     defaultBreak,
		Deadline:  defaultDeadline,
		Debug:     defaultDebug,
		UDPPacket: []byte(defaultUDPPacket),
		Context:   context.Background(),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func WithProto(proto string) Option {
	return func(h *Executor) {
		h.Proto = proto
	}
}

func WithContext(ctx context.Context) Option {
	return func(h *Executor) {
		h.Context = ctx
	}
}
func WithWait(wait time.Duration) Option {
	return func(h *Executor) {
		h.Wait = wait
	}
}

func WithBreak(b time.Duration) Option {
	return func(h *Executor) {
		h.Break = b
	}
}

func WithDeadline(deadline time.Duration) Option {
	return func(h *Executor) {
		h.Deadline = deadline
	}
}

func WithDebug(debug bool) Option {
	return func(h *Executor) {
		h.Debug = debug
	}
}

func WithUDPPacket(packet []byte) Option {
	return func(h *Executor) {
		h.UDPPacket = packet
	}
}

func (e *Executor) Do(addrs []string) bool {
	deadlineCh := time.After(e.Deadline)
	successCh := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(len(addrs))

	go func() {
		for _, addr := range addrs {
			go func(addr string) {
				defer wg.Done()

				for {
					select {
					case <-deadlineCh:
						return
					case <-e.Context.Done():
						return
					default:
						if e.Proto == "udp" {
							if !e.doUDP(addr) {
								continue
							}
						} else if !e.doTCP(addr) {
							continue
						}

						if e.Debug {
							log.Printf("%s %s is OK", e.Proto, addr)
						}

						return
					}
				}
			}(addr)
		}

		wg.Wait()
		close(successCh)
	}()

	select {
	case <-e.Context.Done():
		return false
	case <-deadlineCh:
		return false
	case <-successCh:
		return true
	}
}

func (e *Executor) doTCP(addr string) bool {
	conn, err := net.DialTimeout(e.Proto, addr, e.Wait)
	if err != nil {
		e.processFail(addr)

		return false
	}

	defer conn.Close()

	return true
}

func (e *Executor) doUDP(addr string) bool {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return false
	}

	conn, err := net.DialTimeout(e.Proto, udpAddr.String(), e.Wait)
	if err != nil {
		e.processFail(addr)

		return false
	}

	defer conn.Close()

	// If UDP packet is set - send it
	if len(e.UDPPacket) > 0 {
		_, err = conn.Write(e.UDPPacket)
		if err != nil {
			e.processFail(addr)

			return false
		}
	}

	// Wait for at least 1 byte response
	d := make([]byte, 1)
	_, err = conn.Read(d)
	if err != nil {
		e.processFail(addr)

		return false
	}

	return true
}

func (e *Executor) processFail(addr string) {
	if e.Debug {
		log.Printf("%s %s is FAILED", e.Proto, addr)
	}

	if e.Break > 0 {
		time.Sleep(e.Break)
	}
}
