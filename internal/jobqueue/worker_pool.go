package jobqueue

import (
	"sync"
	repositories "taskflow/internal/domain/repository"

	"github.com/sirupsen/logrus"
)

type WorkerPool struct {
	numWorkers int
	JobQueue   *JobQueue
	repo       repositories.JobRepository
	wg         sync.WaitGroup
	quit       chan struct{}
}

func NewWorkerPool(numWorkers int, jobQueue *JobQueue, repo repositories.JobRepository) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		JobQueue:   jobQueue,
		repo:       repo,
		quit:       make(chan struct{}),
	}
}

func (wp *WorkerPool) Start() {

	for i := range wp.numWorkers {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.quit)
	wp.wg.Wait()
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	logrus.Println("Worker started:", id) // Add this
	for {
		select {
		case jobID := <-wp.JobQueue.GetQueue(): // Get job ID
			logrus.Println("Worker ", id, "Processing job:", jobID)

			//  Fetch the full job object using the job ID
			job, err := wp.repo.FindByID(jobID)
			if err != nil {
				logrus.Println("Failed to fetch job:", err)
				continue
			}

			//  Execute the job logic
			job.Execute()

			//  Update the job status in the database
			if err := wp.repo.Update(job); err != nil {
				logrus.Println("Failed to update job status:", err)
			}

		case <-wp.quit: //  Graceful shutdown
			logrus.Info("Worker stopping")
			return
		}
	}
}
func (wp *WorkerPool) ActiveWorkers() int {
	return wp.numWorkers
}
