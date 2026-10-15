package main

import (
	"container/heap"
	"time"
)

type Job struct {
	ID       string
	Priority int
	Enqueued time.Time
	sequence uint64
}

type jobHeap struct {
	jobs        []Job
	now         time.Time
	agingPeriod time.Duration
	maxPriority int
}

func (h jobHeap) Len() int { return len(h.jobs) }
func (h jobHeap) Less(i, j int) bool {
	a, b := h.effective(h.jobs[i]), h.effective(h.jobs[j])
	if a == b {
		return h.jobs[i].sequence < h.jobs[j].sequence
	}
	return a > b
}
func (h jobHeap) Swap(i, j int)   { h.jobs[i], h.jobs[j] = h.jobs[j], h.jobs[i] }
func (h *jobHeap) Push(value any) { h.jobs = append(h.jobs, value.(Job)) }
func (h *jobHeap) Pop() any {
	last := len(h.jobs) - 1
	job := h.jobs[last]
	h.jobs = h.jobs[:last]
	return job
}
func (h jobHeap) effective(job Job) int {
	if h.agingPeriod <= 0 {
		return job.Priority
	}
	score := job.Priority + int(h.now.Sub(job.Enqueued)/h.agingPeriod)
	if h.maxPriority > 0 && score > h.maxPriority {
		score = h.maxPriority
	}
	return score
}

type Queue struct {
	h    jobHeap
	next uint64
}

func NewQueue(agingPeriod time.Duration, maxPriority int) *Queue {
	return &Queue{h: jobHeap{agingPeriod: agingPeriod, maxPriority: maxPriority}}
}

func (q *Queue) Enqueue(job Job) {
	q.next++
	job.sequence = q.next
	heap.Push(&q.h, job)
}

func (q *Queue) Pop(now time.Time) (Job, bool) {
	if q.h.Len() == 0 {
		return Job{}, false
	}
	q.h.now = now
	heap.Init(&q.h)
	return heap.Pop(&q.h).(Job), true
}

func (q *Queue) Len() int { return q.h.Len() }
