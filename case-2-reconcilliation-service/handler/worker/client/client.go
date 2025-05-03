package workerclient

import (
	"log"

	"github.com/hibiken/asynq"
)

type WorkerClient interface {
	Enqueue(taskType string, payload []byte, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

type workerClient struct {
	client *asynq.Client
}

func NewWorkerClient(client *asynq.Client) WorkerClient {
	return &workerClient{client}
}

func (c *workerClient) Enqueue(taskType string, payload []byte, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	log.Printf("Enqueueing task: %s", taskType)

	task := asynq.NewTask(taskType, payload)
	taskInfo, err := c.client.Enqueue(task, opts...)
	if err != nil {
		return nil, err
	}

	return taskInfo, nil
}
