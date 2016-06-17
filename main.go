package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	timeout bool
	receive int32 = 0
	pushWG  sync.WaitGroup
	pullWG  sync.WaitGroup
)

type Message struct {
	text string
}

func pullMessage(queue chan Message) {
	defer pullWG.Done()
	for _ = range queue {
		atomic.AddInt32(&receive, 1)
	}
}

func pushMessage(queue chan Message) {
	var msg Message
	defer pushWG.Done()
	for !timeout {
		queue <- msg
	}
}

func main() {
	push := flag.Int("push", 1, "Push queue")
	pull := flag.Int("pull", 1, "Push queue")
	duration := flag.String("dur", "5s", "Duration")
	flag.Parse()

	queue := make(chan Message, 1000000000)

	n := time.Now()
	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println("Queue Size: ", len(queue))
		}
	}()
	go func() {
		d, err := time.ParseDuration(*duration)
		if err != nil {
			log.Panicln("Unable to parse duration")
		}
		time.Sleep(d)
		t := time.Since(n)
		fmt.Println("Duration:", t)
		timeout = true
	}()

	for n := 0; n < *pull; n++ {
		pullWG.Add(1)
		go pullMessage(queue)
	}

	for n := 0; n < *push; n++ {
		pushWG.Add(1)
		go pushMessage(queue)
	}

	pushWG.Wait()
	close(queue)

	pullWG.Wait()

	fmt.Println("Messages: ", receive)
}
