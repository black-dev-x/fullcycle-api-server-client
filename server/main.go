package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"black.dev.x/database"
	"gorm.io/gorm"
)

type Response struct {
	Bid string `json:"bid"`
}
type ResponseExternalApi struct {
	Data DollarPrice `json:"USDBRL"`
}
type DollarPrice struct {
	gorm.Model
	Code               string `json:"code"`
	Codein             string `json:"codein"`
	Name               string `json:"name"`
	High               string `json:"high"`
	Low                string `json:"low"`
	VarBid             string `json:"varBid"`
	PctChange          string `json:"pctChange"`
	Bid                string `json:"bid"`
	Ask                string `json:"ask"`
	ExternalTimestamp  string `json:"timestamp"`
	ExternalCreateDate string `json:"create_date"`
}

func main() {
	migrate()
	http.HandleFunc("GET /", getDollarPrice)
	http.ListenAndServe(":8080", nil)
}

func getDollarPrice(w http.ResponseWriter, r *http.Request) {
	ctxApi := context.Background()
	ctxApi, cancel := context.WithTimeout(ctxApi, time.Millisecond*200)
	defer cancel()
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	req, _ := http.NewRequestWithContext(ctxApi, "GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctxApi.Err() != nil {
			println("Timeout error")
		}
		json.NewEncoder(w).Encode(retrieveLatestDollarPrice())
		return
	}
	response := ResponseExternalApi{}
	error := json.NewDecoder(res.Body).Decode(&response)
	if error != nil {
		json.NewEncoder(w).Encode(retrieveLatestDollarPrice())
		return
	}
	price := response.Data
	saveDollarPrice(&price)
	Response := Response{Bid: price.Bid}
	json.NewEncoder(w).Encode(Response)

}

func retrieveLatestDollarPrice() *Response {
	price := DollarPrice{}
	Response := Response{}
	database.DB.Order("created_at desc").First(&price)
	Response.Bid = price.Bid
	return &Response
}

func saveDollarPrice(price *DollarPrice) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancel()
	error := database.DB.WithContext(ctx).Create(&price).Error
	if error != nil {
		if ctx.Err() != nil {
			println("Timeout error at saving dollar price")
		}
		println("Error saving dollar price")
	}
}

func migrate() {
	database.Load()
	database.DB.AutoMigrate(&DollarPrice{})
}
