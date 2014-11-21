package main

import (
	"fmt"
	"github.com/stefankopieczek/gossip/log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

var (
	// Caller parameters
	caller = &endpoint{
		displayName: "Ryan",
		username:    "ryan",
		host:        "172.18.115.70",
		port:        5070,
		proxy:       "10.26.1.41",
		proxy_port:  6000,
		transport:   "UDP",
	}

	// Callee parameters
	callee = &endpoint{
		displayName: "Ryan's PC",
		username:    "stefan",
		host:        "172.18.115.70",
		port:        5060,
		proxy:       "10.26.1.21",
		proxy_port:  6010,
		transport:   "UDP",
	}
)

func main() {
	var err error
	log.SetDefaultLogLevel(log.SEVERE)
	if err != nil {
		panic(err)
	}

	err = caller.Start()
	if err != nil {
		panic(err)
	}
	err = callee.Start()
	if err != nil {
		panic(err)
	}

	f, err := os.Create("cpuprof.out")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	sub_num := 91000
	caller.dialog.cseq = 1
	successes := 0
	tries := 0
	success := make(chan struct{}, 10)
	ticker := time.NewTicker(time.Second)
	go func() {
		for _ = range success {
			successes++
		}
	}()
	end := time.After(20 * time.Second)
	stop := time.After(10 * time.Second)
	try := time.NewTicker(500 * time.Microsecond)
	go callee.ServeNonInvite()
loop:
	for {
		select {
		case <-try.C:
			tries++
			sub_num++
			caller.dialog.callId = fmt.Sprintf("callid%v", sub_num)
			caller.username = fmt.Sprintf("user%v", sub_num)
			caller.dialog.cseq = uint32(sub_num)
			go func() {
				err := caller.Register(callee)
				if err == nil {
					success <- struct{}{}
				} else {
					//fmt.Printf("Registration failed: %v\n", err.Error())
				}
			}()
		case <-ticker.C:
			fmt.Printf("%v tries. %v successes.\n", tries, successes)
			fmt.Printf("Currently %v goroutines exist.", runtime.NumGoroutine())
		case <-stop:
			try.Stop()
		case <-end:
			f, _ := os.Create("memprof.out")
			pprof.WriteHeapProfile(f)
			f.Close()
			break loop
		default:
		}
	}
}
