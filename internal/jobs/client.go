package jobs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/stanleyh24/vendix/internal/config"

	"github.com/hibiken/asynq"
)

// Client wraps asynq.Client for job enqueueing
type Client struct {
	client *asynq.Client
}

// NewClient creates a new job client
func NewClient(cfg *config.Config) *Client {
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}

	return &Client{
		client: asynq.NewClient(redisOpt),
	}
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.client.Close()
}

// EnqueueSendEmail enqueues an email sending task
func (c *Client) EnqueueSendEmail(ctx context.Context, to, subject, body string) error {
	payload, err := json.Marshal(SendEmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
	})
	if err != nil {
		return err
	}

	task := asynq.NewTask(TaskSendEmail, payload)
	_, err = c.client.EnqueueContext(ctx, task, asynq.Queue("default"))
	return err
}

// EnqueueSignInvoice enqueues an invoice signing task
func (c *Client) EnqueueSignInvoice(ctx context.Context, invoiceID, schema string, delay time.Duration) error {
	payload, err := json.Marshal(SignInvoicePayload{
		InvoiceID: invoiceID,
		Schema:    schema,
	})
	if err != nil {
		return err
	}

	task := asynq.NewTask(TaskSignInvoice, payload)
	_, err = c.client.EnqueueContext(ctx, task,
		asynq.Queue("critical"),
		asynq.ProcessIn(delay),
	)
	return err
}

// EnqueueSendInvoiceToDGII enqueues a task to send invoice to DGII
func (c *Client) EnqueueSendInvoiceToDGII(ctx context.Context, invoiceID, schema string) error {
	payload, err := json.Marshal(SendInvoiceToDGIIPayload{
		InvoiceID: invoiceID,
		Schema:    schema,
	})
	if err != nil {
		return err
	}

	task := asynq.NewTask(TaskSendInvoiceToDGII, payload)
	_, err = c.client.EnqueueContext(ctx, task, asynq.Queue("critical"))
	return err
}
