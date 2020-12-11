package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/TechnoServsEIP/Backend/utils"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

type createCheckoutSessionResponse struct {
	SessionID string `json:"id"`
}

func PaymentNew(w http.ResponseWriter, r *http.Request) {
	fmt.Println("initialing new payment")
	defer r.Body.Close()

	req := struct {
		Email   string `json:"email"`
		Product string `json:"product"`
	}{}
	msgFailure := utils.Message(false, "request failed")
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.Respond(w, msgFailure, 400)
		return
	}

	priceToPaid := int64(100) //first payment offer

	domain := "https://blissful-lamarr-d0eb92.netlify.app/#/checkout"
	params := &stripe.CheckoutSessionParams{
		CustomerEmail: stripe.String(req.Email),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyEUR)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Minecraft Subscription"),
					},
					UnitAmount: stripe.Int64(priceToPaid),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "?success=true"),
		CancelURL:  stripe.String(domain + "?canceled=true"),
	}

	sessionPayment, err := session.New(params)
	if err != nil {
		log.Printf("session.New error: %v", err)
		return
	}

	data := createCheckoutSessionResponse{
		SessionID: sessionPayment.ID,
	}

	fmt.Println("session id: ", data.SessionID)
	js, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func PaymentRenew(w http.ResponseWriter, r *http.Request) {
	fmt.Println("initialing new payment")
	defer r.Body.Close()
	userId := r.Context().Value("user").(uint)

	req := struct {
		Email   string `json:"email"`
		Product string `json:"product"`
	}{}
	msgFailure := utils.Message(false, "request failed")
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.Respond(w, msgFailure, 400)
		return
	}

	priceToPaid := int64(GetTotalToPaidPerMonthByUser(userId))
	fmt.Println("the user have to pay " + strconv.FormatInt(priceToPaid, 10))

	domain := "https://blissful-lamarr-d0eb92.netlify.app/#/checkout"
	params := &stripe.CheckoutSessionParams{
		CustomerEmail: stripe.String(req.Email),
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
					UnitAmount: stripe.Int64(priceToPaid),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "?success=true"),
		CancelURL:  stripe.String(domain + "?canceled=true"),
	}

	sessionPayment, err := session.New(params)
	if err != nil {
		log.Printf("session.New: %v", err)
	}

	data := createCheckoutSessionResponse{
		SessionID: sessionPayment.ID,
	}

	fmt.Println("session id: ", data.SessionID)
	js, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
