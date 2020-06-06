package worker

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/nedrocks/delphisbe/graph/model"

	"github.com/nedrocks/delphisbe/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/sirupsen/logrus"

	"github.com/nedrocks/delphisbe/internal/backend"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Worker struct {
	client     *sqs.SQS
	backend    backend.DelphisBackend
	url        string
	maxWorkers int
}

type Message struct {
	Content model.ImportedContentInput
	Done    func() error
}

type MessageIter interface {
	Next(*Message) bool
	Close() error
}

func NewDripWorker(conf config.Config, backend backend.DelphisBackend, sess *session.Session) *Worker {
	return &Worker{
		client:     sqs.New(sess),
		backend:    backend,
		url:        conf.SQSConfig.DripURL,
		maxWorkers: conf.SQSConfig.MaxWorkers,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	iter := w.ReadMessages(ctx)
	queue := make(chan Message, w.maxWorkers)
	workers := sync.WaitGroup{}
	workers.Add(w.maxWorkers)

	for i := 0; i < w.maxWorkers; i++ {
		logrus.Debugf("starting workers")

		go func() {
			defer workers.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-queue:
					if !ok {
						return
					}
					if err := w.handleMessage(ctx, val); err != nil {
						logrus.WithError(err).Error("failed to handled message")
					}
				}
			}
		}()
	}

	var msg Message
	for iter.Next(&msg) {
		queue <- msg
	}

	close(queue)
	workers.Wait()

	if err := iter.Close(); err != nil {
		return err
	}

	return nil
}

func (w *Worker) ReadMessages(ctx context.Context) MessageIter {
	msgs := make(chan Message, 10)
	errs := make(chan error, 1)

	go func() {
		defer close(msgs)

		for ctx.Err() == nil {
			input := sqs.ReceiveMessageInput{
				MaxNumberOfMessages: aws.Int64(10),
				QueueUrl:            aws.String(w.url),
				VisibilityTimeout:   aws.Int64(20),
				WaitTimeSeconds:     aws.Int64(5),
			}

			resp, err := w.client.ReceiveMessageWithContext(ctx, &input)
			if err != nil {
				errs <- err
				return
			}

			for _, msg := range resp.Messages {
				var tempContent model.ImportedContentInput

				logrus.Debugf("Json Body: %+v\n", *msg.Body)

				if err := json.Unmarshal([]byte(*msg.Body), &tempContent); err != nil {
					logrus.WithError(err).Error("error unmarshaling body")
				}

				logrus.Debugf("Temp: %+v\n", tempContent)

				handler := *msg.ReceiptHandle
				msgs <- Message{
					Content: tempContent,
					Done: func() error {
						deleteInput := sqs.DeleteMessageInput{
							QueueUrl:      aws.String(w.url),
							ReceiptHandle: aws.String(handler),
						}

						_, err := w.client.DeleteMessageWithContext(ctx, &deleteInput)
						return err
					},
				}
			}
		}
	}()

	return &sqsIter{
		msgs: msgs,
		errs: errs,
	}
}

func (w *Worker) handleMessage(ctx context.Context, msg Message) error {
	_, err := w.backend.PutImportedContentAndTags(ctx, msg.Content)
	if err != nil {
		return err
	}

	return msg.Done()
}

type sqsIter struct {
	msgs chan Message
	errs chan error
	err  error
}

func (iter *sqsIter) Next(msg *Message) bool {
	if iter.err != nil {
		return false
	}

	select {
	case val, ok := <-iter.msgs:
		if !ok {
			return false
		}

		*msg = val
		return true
	case err := <-iter.errs:
		iter.err = err
		return false
	}
}

func (iter *sqsIter) Close() error {
	return iter.err
}
