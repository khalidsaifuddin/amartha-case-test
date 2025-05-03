package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/config"
)

type WorkerServer struct {
	cfg       config.Config
	server    *asynq.Server
	scheduler *asynq.Scheduler
}

func NewWorkerServer(redisConf asynq.RedisClientOpt, cfg config.Config) *WorkerServer {
	workerConf := asynq.Config{
		Concurrency: cfg.WorkerConcurrency,
		Queues: map[string]int{
			"default": 3,
		},
	}

	server := asynq.NewServer(redisConf, workerConf)

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}

	scheduler := asynq.NewScheduler(redisConf, &asynq.SchedulerOpts{
		Location: loc,
		EnqueueErrorHandler: func(task *asynq.Task, opts []asynq.Option, err error) {
			log.Printf("scheduler failed to run")
		},
	})

	return &WorkerServer{
		cfg:       cfg,
		server:    server,
		scheduler: scheduler,
	}
}

func (w *WorkerServer) Run() {
	log.Printf("worker server running")

	mux := asynq.NewServeMux()

	if err := w.server.Run(mux); err != nil {
		log.Printf("could not run worker server: %v", err)
	}
}

func (w *WorkerServer) RunScheduler() {
	log.Printf("worker scheduler running")

	if err := w.scheduler.Run(); err != nil {
		log.Printf("could not run scheduler server: %v", err)
	}
}

func InitRedisClientWorker(cfg config.Config) asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr: fmt.Sprintf(`%s:%s`, cfg.RedisHost, cfg.RedisPort),
	}
}
