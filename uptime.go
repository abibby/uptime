package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/abibby/uptime/config"
	"github.com/go-ping/ping"
)

var addrs = []string{
	"www.google.com",
	"amazon.com",
}

func uptime(t time.Time) error {
	_, s := checkAll(addrs)

	f, err := OpenFile(config.LogPath, fmt.Sprintf("%s.csv", t.Format("2006-01-02")))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s\t%d\n", t.Format("2006-01-02T15:04:05"), s)
	if err != nil {
		return err
	}

	return nil
}

func checkAll(addrs []string) (map[string]bool, byte) {
	m := map[string]bool{}
	mtx := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	succeeded := byte(0)

	for _, addr := range addrs {
		wg.Add(1)

		go func(addr string) {
			defer wg.Done()

			succ := check(addr)

			mtx.Lock()
			defer mtx.Unlock()
			if succ {
				succeeded++
			}
			m[addr] = succ
		}(addr)
	}

	wg.Wait()

	return m, succeeded
}

func check(addr string) bool {
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		log.Println(err)
		return false
	}

	pinger.SetPrivileged(true)
	pinger.Count = 1

	err = pinger.Run()
	if err != nil {
		log.Println(err)
		return false
	}

	stats := pinger.Statistics()

	return stats.PacketsSent == stats.PacketsRecv
}
