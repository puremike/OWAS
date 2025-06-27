package payments

import (
	"net/http"
	"time"

	"github.com/stripe/stripe-go/v82"
)

type StripePayment struct {
	StripeSecretKey string
	CancelURL       string
	SuccessURL      string
}

func NewStripePayment(stripeSecretKey, cancelURL, successURL string) *StripePayment {
	stripe.Key = stripeSecretKey
	stripe.SetHTTPClient(&http.Client{
		Timeout: 10 * time.Second,
	})
	return &StripePayment{
		StripeSecretKey: stripeSecretKey,
		CancelURL:       cancelURL,
		SuccessURL:      successURL,
	}
}
