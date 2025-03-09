package bootstrap

import (
	repositories "taskflow/internal/domain/repository"
	"taskflow/internal/jobqueue"
)

func InitWorkerPool(repo repositories.JobRepository) (*jobqueue.JobQueue, *jobqueue.WorkerPool) {
	jobQueue := jobqueue.NewJobQueue(10)
	workerPool := jobqueue.NewWorkerPool(2, jobQueue, repo)
	workerPool.Start()
	return jobQueue, workerPool
}
