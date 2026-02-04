package stripe

import (
	"encoding/json"
	"fmt"

	"github.com/stanleyh24/vendix/internal/logger"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

// HandleWebhook processes Stripe webhook events
func HandleWebhook(payload []byte, signatureHeader, webhookSecret string) (*stripe.Event, error) {
	event, err := webhook.ConstructEvent(payload, signatureHeader, webhookSecret)
	if err != nil {
		return nil, fmt.Errorf("webhook signature verification failed: %w", err)
	}

	return &event, nil
}

// ProcessEvent processes a Stripe event
func ProcessEvent(event *stripe.Event) error {
	switch event.Type {
	case "payment_intent.succeeded":
		return handlePaymentIntentSucceeded(event)
	case "payment_intent.payment_failed":
		return handlePaymentIntentFailed(event)
	case "customer.subscription.created":
		return handleSubscriptionCreated(event)
	case "customer.subscription.updated":
		return handleSubscriptionUpdated(event)
	case "customer.subscription.deleted":
		return handleSubscriptionDeleted(event)
	case "invoice.payment_succeeded":
		return handleInvoicePaymentSucceeded(event)
	case "invoice.payment_failed":
		return handleInvoicePaymentFailed(event)
	default:
		logger.Info("Unhandled event type", "type", event.Type)
	}

	return nil
}

func handlePaymentIntentSucceeded(event *stripe.Event) error {
	var paymentIntent stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
		return err
	}

	logger.Info("Payment intent succeeded",
		"payment_intent_id", paymentIntent.ID,
		"amount", paymentIntent.Amount,
		"customer", paymentIntent.Customer,
	)

	// TODO: Update payment record in database

	return nil
}

func handlePaymentIntentFailed(event *stripe.Event) error {
	var paymentIntent stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
		return err
	}

	logger.Error("Payment intent failed",
		"payment_intent_id", paymentIntent.ID,
		"amount", paymentIntent.Amount,
		"error", paymentIntent.LastPaymentError,
	)

	// TODO: Update payment record and notify customer

	return nil
}

func handleSubscriptionCreated(event *stripe.Event) error {
	var subscription stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		return err
	}

	logger.Info("Subscription created",
		"subscription_id", subscription.ID,
		"customer", subscription.Customer.ID,
	)

	// TODO: Update subscription record in database

	return nil
}

func handleSubscriptionUpdated(event *stripe.Event) error {
	var subscription stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		return err
	}

	logger.Info("Subscription updated",
		"subscription_id", subscription.ID,
		"status", subscription.Status,
	)

	// TODO: Update subscription record in database

	return nil
}

func handleSubscriptionDeleted(event *stripe.Event) error {
	var subscription stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		return err
	}

	logger.Info("Subscription deleted",
		"subscription_id", subscription.ID,
	)

	// TODO: Update subscription status in database

	return nil
}

func handleInvoicePaymentSucceeded(event *stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return err
	}

	logger.Info("Invoice payment succeeded",
		"invoice_id", invoice.ID,
		"subscription", invoice.Subscription,
	)

	// TODO: Record payment and update subscription

	return nil
}

func handleInvoicePaymentFailed(event *stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return err
	}

	logger.Error("Invoice payment failed",
		"invoice_id", invoice.ID,
		"subscription", invoice.Subscription,
	)

	// TODO: Notify customer of payment failure

	return nil
}
