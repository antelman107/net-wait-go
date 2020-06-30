package wait

import (
	"log"
	"net"
	"sync"
	"time"
)

func Do(proto string, addrs []string, delay, deadline time.Duration, debug bool) bool {
	deadlineCh := time.After(deadline)
	successCh := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(len(addrs))

	go func() {
		for _, addr := range addrs {
			go func(addr string) {
				defer wg.Done()

				for {
					conn, err := net.DialTimeout(proto, addr, delay)
					if err != nil {
						if debug {
							log.Printf("%s is FAILED", addr)
						}

						time.Sleep(delay)

						continue
					}

					if debug {
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
