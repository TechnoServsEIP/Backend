package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"log"
	"net/http"
)

type createCheckoutSessionResponse struct {
	SessionID string `json:"id"`
}

var PaymentNew = func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("initialing new payment")
	domain := "https://technoservs.ichbinkour.eu/#/checkout" //TODO change this
	params := &stripe.CheckoutSessionParams{
		CustomerEmail: stripe.String("jonathan.frickert@epitech.eu"),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyEUR)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Minecraft"),
					},
					UnitAmount: stripe.Int64(1200),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "?success=true"),
		CancelURL:  stripe.String(domain + "?canceled=true"),
	}
	session, err := session.New(params)
	if err != nil {
		log.Printf("session.New: %v", err)
	}
	data := createCheckoutSessionResponse{
		SessionID: session.ID,
	}
	fmt.Println("session id: ", data.SessionID)
	js, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
