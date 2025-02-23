package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Bid string `json:"bid"`
}

func main() {
	ctxApi := context.Background()
	ctxApi, cancel := context.WithTimeout(ctxApi, time.Millisecond*300)
	req, error := http.NewRequestWithContext(ctxApi, "GET", "http://localhost:8080/", nil)
	if error != nil {
		panic(error)
	}
	defer cancel()
	res, error := http.DefaultClient.Do(req)
	if error != nil {
		panic(error)
	}
	defer res.Body.Close()
	response := Response{}
	err := json.NewDecoder(res.Body).Decode(&response)
	println(response.Bid)
	if err != nil {
		panic(err)
	}
	text := "Dólar: " + response.Bid
	file, error := os.Create("cotacao.txt")
	if error != nil {
		panic(error)
	}
	defer file.Close()
	file.WriteString(text)
}
