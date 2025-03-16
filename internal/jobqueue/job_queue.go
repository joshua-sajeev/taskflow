package jobqueue

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type JobQueue struct {
	queue chan uuid.UUID
}

func NewJobQueue(size int) *JobQueue {
	return &JobQueue{
		queue: make(chan uuid.UUID, size),
	}
}

func (jq *JobQueue) Enqueue(jobID uuid.UUID) {
	select {
	case jq.queue <- jobID:
		logrus.Println("Job enqueued:", jobID, "Queue size", jq.Size())
	default:
		logrus.Println("Queue full! Dropping job:", jobID)
	}
}

func (jq *JobQueue) Dequeue() uuid.UUID {
	return <-jq.queue
}

func (jq *JobQueue) GetQueue() chan uuid.UUID {
	return jq.queue
}

func (jq *JobQueue) Size() int {
	return len(jq.queue)
}
