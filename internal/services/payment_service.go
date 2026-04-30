package services

import (
	"os"
	"strconv"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

// 🔥 cria checkout do Stripe com invoice_id
func CreateCheckoutSession(amount float64, invoiceID int) (string, error) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),

		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),

					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Invoice Payment"),
					},

					UnitAmount: stripe.Int64(int64(amount * 100)), // centavos
				},
				Quantity: stripe.Int64(1),
			},
		},

		// 🔥 MUITO IMPORTANTE
		Metadata: map[string]string{
			"invoice_id": strconv.Itoa(invoiceID),
		},

		SuccessURL: stripe.String("http://localhost:3000/success"),
		CancelURL:  stripe.String("http://localhost:3000/cancel"),
	}

	s, err := session.New(params)
	if err != nil {
		return "", err
	}

	return s.URL, nil
}