package main

import (
	"testing"
	"time"
)

func TestAgingLetsOldLowPriorityJobRun(t *testing.T) {
	start := time.Unix(0, 0)
	q := NewQueue(time.Second, 10)
	q.Enqueue(Job{ID: "old", Priority: 1, Enqueued: start})
	q.Enqueue(Job{ID: "new", Priority: 10, Enqueued: start.Add(20 * time.Second)})
	got, ok := q.Pop(start.Add(20 * time.Second))
	if !ok || got.ID != "old" {
		t.Fatalf("got=%+v ok=%t", got, ok)
	}
}

func TestEqualPriorityIsFIFO(t *testing.T) {
	now := time.Now()
	q := NewQueue(0, 0)
	q.Enqueue(Job{ID: "first", Priority: 3, Enqueued: now})
	q.Enqueue(Job{ID: "second", Priority: 3, Enqueued: now})
	got, _ := q.Pop(now)
	if got.ID != "first" {
		t.Fatalf("got=%s", got.ID)
	}
}
