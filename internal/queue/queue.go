package queue

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rzhaka-turiki/apex_account_verifier/internal/api"
	"github.com/rzhaka-turiki/apex_account_verifier/internal/model"
	"golang.org/x/time/rate"
)

var ErrQueueFull = errors.New("request queue is full")

type Job struct {
	ctx    context.Context
	fn     func(context.Context) (*model.Account, error)
	result chan Result
}

type Result struct {
	account *model.Account
	err     error
}

type Queue struct {
	jobs    chan Job
	limiter *rate.Limiter
}

func NewQueue(capacity, workers, requestsPerSecond int) *Queue {
	q := &Queue{
		jobs:    make(chan Job, capacity),
		limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), 1),
	}
	for i := 0; i < workers; i++ {
		go q.worker()
	}
	return q
}

func (q *Queue) Submit(ctx context.Context, fn func(context.Context) (*model.Account, error)) (*model.Account, error) {
	job := Job{
		ctx:    ctx,
		fn:     fn,
		result: make(chan Result, 1),
	}
	select {
	case q.jobs <- job:
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return nil, ErrQueueFull
	}
	select {
	case result := <-job.result:
		return result.account, result.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (q *Queue) execute(job Job) (*model.Account, error) {
	const maxRetries = 3
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if err := q.limiter.Wait(job.ctx); err != nil {
			return nil, err
		}
		account, err := job.fn(job.ctx)
		if err == nil {
			return account, nil
		}
		var apiErr *api.APIError
		if !errors.As(err, &apiErr) {
			return nil, err
		}
		if apiErr.StatusCode != 429 {
			return nil, err
		}
		if attempt == maxRetries {
			return nil, err
		}
		delay := time.Duration(1<<attempt) * time.Second
		fmt.Printf("rate limit exceeded, retrying in %s (attempt %d/%d)\n",
			delay,
			attempt+1,
			maxRetries)
		select {
		case <-time.After(delay):
			continue
		case <-job.ctx.Done():
			return nil, job.ctx.Err()
		}
	}
	return nil, errors.New("unexpected retry loop exit")
}

func (q *Queue) worker() {
	for job := range q.jobs {
		account, err := q.execute(job)
		job.result <- Result{
			account: account,
			err:     err,
		}
	}
}
