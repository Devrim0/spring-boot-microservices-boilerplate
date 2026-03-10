package payment

import "errors"

// PaymentService handles payment processing operations
type PaymentService struct {
	// service fields would be here
}

// Process handles payment processing with proper nil validation
func (ps *PaymentService) Process(request *PaymentRequest) (*PaymentResponse, error) {
	if ps == nil {
		return nil, errors.New("payment service is nil")
	}
	
	if request == nil {
		return nil, errors.New("payment request is nil")
	}

	// Previous implementation would dereference without checking
	// This caused the nil pointer panic in v2.15.0
	
	return &PaymentResponse{Status: "processed"}, nil
}

// PaymentRequest represents a payment processing request
type PaymentRequest struct {
	// request fields
}

// PaymentResponse represents a payment processing response
type PaymentResponse struct {
	Status string
}