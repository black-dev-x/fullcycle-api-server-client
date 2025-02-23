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
	Data DolarPrice `json:"USDBRL"`
}
type DolarPrice struct {
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
	http.HandleFunc("GET /", getDolarPrice)
	http.ListenAndServe(":8080", nil)
}

func getDolarPrice(w http.ResponseWriter, r *http.Request) {
	ctxApi := context.Background()
	ctxApi, cancel := context.WithTimeout(ctxApi, time.Millisecond*200)
	defer cancel()
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	req, err := http.NewRequestWithContext(ctxApi, "GET", url, nil)
	if err != nil {
		json.NewEncoder(w).Encode(retrieveLatestDolarPrice())
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		json.NewEncoder(w).Encode(retrieveLatestDolarPrice())
		return
	}
	response := ResponseExternalApi{}
	error := json.NewDecoder(res.Body).Decode(&response)
	if error != nil {
		json.NewEncoder(w).Encode(retrieveLatestDolarPrice())
		return
	}
	price := response.Data
	saveDolarPrice(&price)
	Response := Response{Bid: price.Bid}
	json.NewEncoder(w).Encode(Response)

}

func retrieveLatestDolarPrice() *Response {
	price := DolarPrice{}
	Response := Response{}
	database.DB.Order("created_at desc").First(&price)
	Response.Bid = price.Bid
	return &Response
}

func saveDolarPrice(price *DolarPrice) {
	ctxApi := context.Background()
	ctxApi, cancel := context.WithTimeout(ctxApi, time.Millisecond*200)
	defer cancel()
	database.DB.WithContext(ctxApi).Create(&price)
}

func migrate() {
	database.Load()
	database.DB.AutoMigrate(&DolarPrice{})
}
