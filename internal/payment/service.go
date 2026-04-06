package payment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shop/checkout/internal/config"
	"github.com/shop/checkout/internal/metrics"
	"github.com/shop/checkout/internal/tracing"
	"github.com/shop/checkout/pkg/logger"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

type PaymentService struct {
	config     *config.Config
	logger     logger.Logger
	metrics    *metrics.Client
	stripeKey  string
}

type Payment struct {
	ID          string  `json:"id"`
	Amount      int64   `json:"amount"`
	Currency    string  `json:"currency"`
	CustomerID  string  `json:"customer_id"`
	Description string  `json:"description"`
	Metadata    map[string]string `json:"metadata"`
}

type ProcessResult struct {
	PaymentIntentID string            `json:"payment_intent_id"`
	Status          string            `json:"status"`
	ClientSecret    string            `json:"client_secret"`
	Metadata        map[string]string `json:"metadata"`
	ProcessedAt     time.Time         `json:"processed_at"`
}

func NewPaymentService(cfg *config.Config, log logger.Logger, m *metrics.Client) *PaymentService {
	return &PaymentService{
		config:    cfg,
		logger:    log,
		metrics:   m,
		stripeKey: cfg.StripeSecretKey,
	}
}

// Process handles payment processing for checkout requests
func (p *PaymentService) Process(ctx context.Context, payment *Payment) (*ProcessResult, error) {
	// Add nil check to prevent panic
	if payment == nil {
		return nil, errors.New("payment object cannot be nil")
	}

	span, ctx := tracing.StartSpan(ctx, "payment.process")
	defer span.Finish()

	p.logger.Info("Processing payment", map[string]interface{}{
		"payment_id":  payment.ID,
		"customer_id": payment.CustomerID,
		"amount":      payment.Amount,
		"currency":    payment.Currency,
	})

	// Record metrics
	p.metrics.Inc("payment.process.started", map[string]string{
		"currency": payment.Currency,
	})

	start := time.Now()
	defer func() {
		p.metrics.Histogram("payment.process.duration", time.Since(start).Milliseconds(), map[string]string{
			"currency": payment.Currency,
		})
	}()

	// Validate payment data
	if err := p.validatePayment(payment); err != nil {
		p.logger.Error("Payment validation failed", map[string]interface{}{
			"payment_id": payment.ID,
			"error":      err.Error(),
		})
		p.metrics.Inc("payment.process.validation_failed", map[string]string{
			"currency": payment.Currency,
		})
		return nil, fmt.Errorf("payment validation failed: %w", err)
	}

	// Set Stripe API key
	stripe.Key = p.stripeKey

	// Create payment intent parameters
	params := &stripe.PaymentIntentParams{
		Amount:      stripe.Int64(payment.Amount),
		Currency:    stripe.String(payment.Currency),
		Customer:    stripe.String(payment.CustomerID),
		Description: stripe.String(payment.Description),
	}

	// Add metadata
	if payment.Metadata != nil {
		for key, value := range payment.Metadata {
			params.AddMetadata(key, value)
		}
	}

	// Add internal tracking metadata
	params.AddMetadata("internal_payment_id", payment.ID)
	params.AddMetadata("processed_at", time.Now().UTC().Format(time.RFC3339))

	// Create payment intent with Stripe
	pi, err := paymentintent.New(params)
	if err != nil {
		p.logger.Error("Stripe payment intent creation failed", map[string]interface{}{
			"payment_id": payment.ID,
			"error":      err.Error(),
		})
		p.metrics.Inc("payment.process.stripe_failed", map[string]string{
			"currency": payment.Currency,
		})
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	p.logger.Info("Payment intent created successfully", map[string]interface{}{
		"payment_id":         payment.ID,
		"payment_intent_id":  pi.ID,
		"status":             string(pi.Status),
		"client_secret":      pi.ClientSecret,
	})

	p.metrics.Inc("payment.process.success", map[string]string{
		"currency": payment.Currency,
		"status":   string(pi.Status),
	})

	// Build result
	result := &ProcessResult{
		PaymentIntentID: pi.ID,
		Status:          string(pi.Status),
		ClientSecret:    pi.ClientSecret,
		Metadata:        make(map[string]string),
		ProcessedAt:     time.Now().UTC(),
	}

	// Copy metadata from Stripe response
	for key, value := range pi.Metadata {
		result.Metadata[key] = value
	}

	return result, nil
}

func (p *PaymentService) validatePayment(payment *Payment) error {
	if payment.ID == "" {
		return errors.New("payment ID is required")
	}

	if payment.Amount <= 0 {
		return errors.New("payment amount must be greater than 0")
	}

	if payment.Currency == "" {
		return errors.New("payment currency is required")
	}

	if payment.CustomerID == "" {
		return errors.New("customer ID is required")
	}

	// Validate currency format
	if len(payment.Currency) != 3 {
		return errors.New("currency must be a 3-character ISO code")
	}

	// Validate amount limits (e.g., Stripe minimum)
	if payment.Currency == "usd" && payment.Amount < 50 {
		return errors.New("minimum amount for USD is $0.50")
	}

	return nil
}

// GetPaymentIntent retrieves a payment intent by ID
func (p *PaymentService) GetPaymentIntent(ctx context.Context, paymentIntentID string) (*stripe.PaymentIntent, error) {
	if paymentIntentID == "" {
		return nil, errors.New("payment intent ID is required")
	}

	span, ctx := tracing.StartSpan(ctx, "payment.get_intent")
	defer span.Finish()

	stripe.Key = p.stripeKey

	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		p.logger.Error("Failed to retrieve payment intent", map[string]interface{}{
			"payment_intent_id": paymentIntentID,
			"error":             err.Error(),
		})
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}

	return pi, nil
}

// CancelPaymentIntent cancels a payment intent
func (p *PaymentService) CancelPaymentIntent(ctx context.Context, paymentIntentID string) error {
	if paymentIntentID == "" {
		return errors.New("payment intent ID is required")
	}

	span, ctx := tracing.StartSpan(ctx, "payment.cancel_intent")
	defer span.Finish()

	stripe.Key = p.stripeKey

	params := &stripe.PaymentIntentCancelParams{}
	_, err := paymentintent.Cancel(paymentIntentID, params)
	if err != nil {
		p.logger.Error("Failed to cancel payment intent", map[string]interface{}{
			"payment_intent_id": paymentIntentID,
			"error":             err.Error(),
		})
		return fmt.Errorf("failed to cancel payment intent: %w", err)
	}

	p.logger.Info("Payment intent cancelled successfully", map[string]interface{}{
		"payment_intent_id": paymentIntentID,
	})

	return nil
}