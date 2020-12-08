package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/TechnoServsEIP/Backend/models"
	"github.com/TechnoServsEIP/Backend/utils"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"log"
	"net/http"
	"strconv"
	"time"
)

type createCheckoutSessionResponse struct {
	SessionID string `json:"id"`
}

func PaymentNew(w http.ResponseWriter, r *http.Request) {
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

	priceToPaid := int64(6) //first payment offer

	domain := "https://technoservs.ichbinkour.eu/#/checkout"
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
	currentTime := time.Now()
	bill := &models.Bill{
		UserId:       userId,
		Email:        req.Email,
		Price:        strconv.FormatInt(priceToPaid, 10),
		Product:      req.Product + " first time subscription",
		StartSubDate: currentTime,
		EndSubDate:   currentTime.AddDate(0, 1, 0),
	}
	bill.InsertBill()

	fmt.Println(*bill)
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

	domain := "https://technoservs.ichbinkour.eu/#/checkout"
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
	currentTime := time.Now()
	bill := &models.Bill{
		UserId:       userId,
		Email:        req.Email,
		Price:        strconv.FormatInt(priceToPaid, 10),
		Product:      req.Product,
		StartSubDate: currentTime,
		EndSubDate:   currentTime.AddDate(0, 1, 0),
	}
	bill.InsertBill()

	fmt.Println(*bill)
	fmt.Println("session id: ", data.SessionID)
	js, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
