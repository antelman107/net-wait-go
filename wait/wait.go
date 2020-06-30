package wait

import (
	"log"
	"net"
	"sync"
	"time"
)

type Executor struct {
	Proto    string
	Addrs    []string
	Wait     time.Duration
	Break    time.Duration
	Deadline time.Duration
	Debug    bool
}

type Option func(*Executor)

func New(opts ...Option) *Executor {
	const (
		defaultProto    = "tcp"
		defaultWait     = 200 * time.Millisecond
		defaultBreak    = 50 * time.Millisecond
		defaultDeadline = 15 * time.Second
		defaultDebug    = false
	)

	e := &Executor{
		Proto:    defaultProto,
		Wait:     defaultWait,
		Break:    defaultBreak,
		Deadline: defaultDeadline,
		Debug:    defaultDebug,
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
					conn, err := net.DialTimeout(e.Proto, addr, e.Wait)
					if err != nil {
						if e.Debug {
							log.Printf("%s is FAILED", addr)
						}

						if e.Break > 0 {
							time.Sleep(e.Break)
						}

						continue
					}

					if e.Debug {
						log.Printf("%s is OK", addr)
					}

					_ = conn.Close()

					return
				}
			}(addr)
		}

		wg.Wait()
		close(successCh)
	}()

	select {
	case <-deadlineCh:
		return false
	case <-successCh:
		return true
	}
}
