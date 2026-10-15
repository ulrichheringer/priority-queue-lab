package main

import (
	"flag"
	"fmt"
	"time"
)

type arrival struct {
	job Job
	at  time.Duration
}

func main() {
	aging := flag.Duration("aging", 0, "priority gain interval; zero disables aging")
	flag.Parse()
	simulate(*aging)
}

func simulate(aging time.Duration) {
	const maxPriority = 10
	var arrivals []arrival
	arrivals = append(arrivals, arrival{Job{ID: "low", Priority: 1}, 0})
	for i := 0; i < 12; i++ {
		arrivals = append(arrivals, arrival{Job{ID: fmt.Sprintf("high-%02d", i+1), Priority: 10}, time.Duration(i) * time.Second})
	}
	base := time.Unix(0, 0)
	queue := NewQueue(aging, maxPriority)
	now := time.Duration(0)
	next, service := 0, 2*time.Second
	fmt.Printf("aging=%s service=%s\n", aging, service)
	for next < len(arrivals) || queue.Len() > 0 {
		for next < len(arrivals) && arrivals[next].at <= now {
			a := arrivals[next]
			a.job.Enqueued = base.Add(a.at)
			queue.Enqueue(a.job)
			next++
		}
		if queue.Len() == 0 {
			now = arrivals[next].at
			continue
		}
		started := base.Add(now)
		job, _ := queue.Pop(started)
		fmt.Printf("start=%2s job=%-8s priority=%2d wait=%2s\n", now, job.ID, job.Priority, started.Sub(job.Enqueued))
		now += service
	}
}
