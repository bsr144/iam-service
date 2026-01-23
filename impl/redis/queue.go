package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func NewJob(payload any, maxRetry int) (*Job, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return &Job{
		ID:        uuid.New().String(),
		Payload:   data,
		Attempts:  0,
		MaxRetry:  maxRetry,
		CreatedAt: time.Now(),
	}, nil
}

func (r *Redis) Enqueue(ctx context.Context, queueName string, job *Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return ErrFailedToMarshalValue(err)
	}
	return r.client.RPush(ctx, queueKey(queueName), data).Err()
}

func (r *Redis) EnqueuePayload(ctx context.Context, queueName string, payload any, maxRetry int) (string, error) {
	job, err := NewJob(payload, maxRetry)
	if err != nil {
		return "", err
	}
	if err := r.Enqueue(ctx, queueName, job); err != nil {
		return "", err
	}
	return job.ID, nil
}

func (r *Redis) EnqueueDelayed(ctx context.Context, queueName string, job *Job, delay time.Duration) error {
	data, err := json.Marshal(job)
	if err != nil {
		return ErrFailedToMarshalValue(err)
	}
	score := float64(time.Now().Add(delay).Unix())
	return r.client.ZAdd(ctx, delayedKey(queueName), redis.Z{
		Score:  score,
		Member: data,
	}).Err()
}

func (r *Redis) Dequeue(ctx context.Context, queueName string, timeout time.Duration) (*Job, error) {

	data, err := r.client.BRPopLPush(ctx, queueKey(queueName), processingKey(queueName), timeout).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrQueueEmpty
		}
		return nil, fmt.Errorf("failed to dequeue: %w", err)
	}

	var job Job
	if err := json.Unmarshal(data, &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	job.Attempts++
	return &job, nil
}

func (r *Redis) DequeueBlocking(ctx context.Context, queueName string) (*Job, error) {
	return r.Dequeue(ctx, queueName, 0)
}

func (r *Redis) AcknowledgeJob(ctx context.Context, queueName string, job *Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	job.Attempts--
	origData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal original job: %w", err)
	}

	removed, err := r.client.LRem(ctx, processingKey(queueName), 1, origData).Result()
	if err != nil {
		return fmt.Errorf("failed to acknowledge job: %w", err)
	}
	if removed == 0 {

		_, err = r.client.LRem(ctx, processingKey(queueName), 1, data).Result()
		if err != nil {
			return fmt.Errorf("failed to acknowledge job: %w", err)
		}
	}
	return nil
}

func (r *Redis) RejectJob(ctx context.Context, queueName string, job *Job, errMsg string) error {
	job.Error = errMsg

	origJob := *job
	origJob.Attempts--
	origData, _ := json.Marshal(&origJob)
	r.client.LRem(ctx, processingKey(queueName), 1, origData)

	if job.Attempts >= job.MaxRetry {

		data, err := json.Marshal(job)
		if err != nil {
			return fmt.Errorf("failed to marshal job: %w", err)
		}
		return r.client.RPush(ctx, deadLetterKey(queueName), data).Err()
	}

	return r.Enqueue(ctx, queueName, job)
}

func (r *Redis) ProcessDelayed(ctx context.Context, queueName string) (int, error) {
	now := float64(time.Now().Unix())

	jobs, err := r.client.ZRangeByScore(ctx, delayedKey(queueName), &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%f", now),
	}).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get delayed jobs: %w", err)
	}

	if len(jobs) == 0 {
		return 0, nil
	}

	pipe := r.client.Pipeline()
	for _, jobData := range jobs {
		pipe.RPush(ctx, queueKey(queueName), jobData)
		pipe.ZRem(ctx, delayedKey(queueName), jobData)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to move delayed jobs: %w", err)
	}

	return len(jobs), nil
}

func (r *Redis) QueueLen(ctx context.Context, queueName string) (int64, error) {
	return r.client.LLen(ctx, queueKey(queueName)).Result()
}

func (r *Redis) ProcessingLen(ctx context.Context, queueName string) (int64, error) {
	return r.client.LLen(ctx, processingKey(queueName)).Result()
}

func (r *Redis) DeadLetterLen(ctx context.Context, queueName string) (int64, error) {
	return r.client.LLen(ctx, deadLetterKey(queueName)).Result()
}

func (r *Redis) DelayedLen(ctx context.Context, queueName string) (int64, error) {
	return r.client.ZCard(ctx, delayedKey(queueName)).Result()
}

func (r *Redis) ClearQueue(ctx context.Context, queueName string) error {
	pipe := r.client.Pipeline()
	pipe.Del(ctx, queueKey(queueName))
	pipe.Del(ctx, processingKey(queueName))
	pipe.Del(ctx, delayedKey(queueName))
	_, err := pipe.Exec(ctx)
	return err
}

func (r *Redis) RequeueProcessing(ctx context.Context, queueName string) (int64, error) {
	var count int64
	for {
		result, err := r.client.RPopLPush(ctx, processingKey(queueName), queueKey(queueName)).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				break
			}
			return count, fmt.Errorf("failed to requeue processing job: %w", err)
		}
		if result == "" {
			break
		}
		count++
	}
	return count, nil
}

func (r *Redis) GetDeadLetterJobs(ctx context.Context, queueName string, start, stop int64) ([]*Job, error) {
	data, err := r.client.LRange(ctx, deadLetterKey(queueName), start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get dead letter jobs: %w", err)
	}

	jobs := make([]*Job, 0, len(data))
	for _, d := range data {
		var job Job
		if err := json.Unmarshal([]byte(d), &job); err != nil {
			continue
		}
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

func (r *Redis) RetryDeadLetter(ctx context.Context, queueName string, jobID string) error {

	data, err := r.client.LRange(ctx, deadLetterKey(queueName), 0, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get dead letter jobs: %w", err)
	}

	for _, d := range data {
		var job Job
		if err := json.Unmarshal([]byte(d), &job); err != nil {
			continue
		}
		if job.ID == jobID {

			job.Attempts = 0
			job.Error = ""
			if err := r.Enqueue(ctx, queueName, &job); err != nil {
				return err
			}

			return r.client.LRem(ctx, deadLetterKey(queueName), 1, d).Err()
		}
	}
	return ErrJobNotFound
}

var ErrQueueEmpty = errors.New("queue is empty")

var ErrJobNotFound = errors.New("job not found")
