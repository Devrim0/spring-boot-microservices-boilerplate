package payment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shop/checkout/internal/domain"
	"github.com/shop/checkout/internal/stripe"
	"go.uber.org/zap"
)

type Service struct {
	stripeClient stripe.Client
	logger       *zap.Logger
}

func NewService(stripeClient stripe.Client, logger *zap.Logger) *Service {
	return &Service{
		stripeClient: stripeClient,
		logger:       logger,
	}
}

type ProcessRequest struct {
	Payment *domain.Payment `json:"payment"`
	UserID  string          `json:"user_id"`
}

type ProcessResponse struct {
	TransactionID string                 `json:"transaction_id"`
	Status        domain.PaymentStatus   `json:"status"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

func (p *Service) Process(ctx context.Context, req *ProcessRequest) (*ProcessResponse, error) {
	p.logger.Info("processing payment", "user_id", req.UserID)

	// Validate request
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	if err := p.validatePaymentRequest(req); err != nil {
		return nil, fmt.Errorf("invalid payment request: %w", err)
	}

	// Ensure payment object is properly initialized
	if req.Payment == nil {
		return nil, fmt.Errorf("payment object cannot be nil")
	}
	if req.Payment.Amount <= 0 {
		return nil, fmt.Errorf("payment amount must be positive")
	}

	// Process payment with Stripe
	stripeResp, err := p.stripeClient.ProcessPayment(ctx, &stripe.PaymentRequest{
		Amount:   req.Payment.Amount,
		Currency: req.Payment.Currency,
		Token:    req.Payment.Token,
	})
	if err != nil {
		// Handle timeout scenarios gracefully
		if errors.Is(err, context.DeadlineExceeded) {
			p.logger.Error("payment timeout - will retry", "error", err, "payment_id", req.Payment.ID)
			return nil, fmt.Errorf("payment processing timeout: %w", err)
		}
		
		return nil, fmt.Errorf("stripe payment failed: %w", err)
	}

	// Create response
	resp := &ProcessResponse{
		TransactionID: stripeResp.TransactionID,
		Status:        domain.PaymentStatusCompleted,
		Metadata: map[string]interface{}{
			"stripe_payment_id": stripeResp.PaymentID,
			"processed_at":      time.Now().UTC(),
		},
	}

	p.logger.Info("payment processed successfully", 
		"transaction_id", resp.TransactionID,
		"user_id", req.UserID,
		"amount", req.Payment.Amount)

	return resp, nil
}

func (p *Service) validatePaymentRequest(req *ProcessRequest) error {
	if req.Payment == nil {
		return fmt.Errorf("payment is required")
	}
	
	if req.Payment.Amount <= 0 {
		return fmt.Errorf("payment amount must be positive")
	}
	
	if req.Payment.Currency == "" {
		return fmt.Errorf("payment currency is required")
	}
	
	if req.Payment.Token == "" {
		return fmt.Errorf("payment token is required")
	}
	
	return nil
}