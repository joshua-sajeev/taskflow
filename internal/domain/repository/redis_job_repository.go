package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"taskflow/internal/domain/entities"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RedisJobRepository is the Redis implementation of JobRepository.
type RedisJobRepository struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisJobRepository creates a new instance of RedisJobRepository.
func NewRedisJobRepository(client *redis.Client) JobRepository {
	return &RedisJobRepository{
		client: client,
		ctx:    context.Background(),
	}
}

// Create inserts a new job into Redis.
func (r *RedisJobRepository) Create(job *entities.Job) error {
	job.ID = uuid.New()
	job.CreatedAt = time.Now()

	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	// Store in Redis with ID as key
	key := fmt.Sprintf("job:%s", job.ID.String())
	return r.client.Set(r.ctx, key, data, 0).Err()
}

// FindByID retrieves a job by ID.
func (r *RedisJobRepository) FindByID(id uuid.UUID) (*entities.Job, error) {
	key := fmt.Sprintf("job:%s", id.String())
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var job entities.Job
	if err := json.Unmarshal([]byte(data), &job); err != nil {
		return nil, err
	}

	return &job, nil
}

// Update modifies an existing job.
func (r *RedisJobRepository) Update(job *entities.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("job:%s", job.ID.String())
	return r.client.Set(r.ctx, key, data, 0).Err()
}

// Delete removes a job by ID.
func (r *RedisJobRepository) Delete(id uuid.UUID) error {
	key := fmt.Sprintf("job:%s", id.String())
	return r.client.Del(r.ctx, key).Err()
}

// CountJobsByStatus counts pending and completed jobs.
func (r *RedisJobRepository) CountJobsByStatus() (pending int64, completed int64, err error) {
	keys, err := r.client.Keys(r.ctx, "job:*").Result()
	if err != nil {
		return 0, 0, err
	}

	for _, key := range keys {
		data, err := r.client.Get(r.ctx, key).Result()
		if err != nil {
			return 0, 0, err
		}

		var job entities.Job
		if err := json.Unmarshal([]byte(data), &job); err != nil {
			return 0, 0, err
		}

		if job.Status == "pending" {
			pending++
		} else if job.Status == "completed" {
			completed++
		}
	}

	return pending, completed, nil
}
