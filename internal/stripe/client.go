package stripe

import (
	"vendix/internal/config"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/subscription"
)

// Client wraps Stripe API client
type Client struct {
	cfg *config.Config
}

// NewClient creates a new Stripe client
func NewClient(cfg *config.Config) *Client {
	stripe.Key = cfg.StripeSecretKey
	return &Client{cfg: cfg}
}

// CreateCustomer creates a Stripe customer
func (c *Client) CreateCustomer(email, name string, metadata map[string]string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
	}

	for k, v := range metadata {
		params.AddMetadata(k, v)
	}

	return customer.New(params)
}

// CreatePaymentIntent creates a payment intent
func (c *Client) CreatePaymentIntent(amount int64, currency, customerID string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
		Customer: stripe.String(customerID),
	}

	return paymentintent.New(params)
}

// CreateSubscription creates a subscription
func (c *Client) CreateSubscription(customerID, priceID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
	}

	return subscription.New(params)
}

// CancelSubscription cancels a subscription
func (c *Client) CancelSubscription(subscriptionID string) (*stripe.Subscription, error) {
	return subscription.Cancel(subscriptionID, nil)
}

// GetCustomer retrieves a customer
func (c *Client) GetCustomer(customerID string) (*stripe.Customer, error) {
	return customer.Get(customerID, nil)
}

// GetSubscription retrieves a subscription
func (c *Client) GetSubscription(subscriptionID string) (*stripe.Subscription, error) {
	return subscription.Get(subscriptionID, nil)
}
