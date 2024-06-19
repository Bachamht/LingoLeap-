package main

import (
	"QuantitativeTrading/db"
	"QuantitativeTrading/om"
	"QuantitativeTrading/tg"
	"QuantitativeTrading/trade"
	"log"
	"github.com/joho/godotenv"
	"os"
	"fmt"
	"strconv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("env load error:", err)
	}
			
	apiKey := os.Getenv("apiKey")
	secretKey := os.Getenv("secretKey")
	baseURL := os.Getenv("baseURL")
	User := os.Getenv("User")
	Port := os.Getenv("Port")
	Password := os.Getenv("Password")
	Host := os.Getenv("Host")
	Dbname := os.Getenv("Dbname")
	dbURL := User + ":" + Password + "@tcp(" + Host + ":" + Port + ")/" + Dbname
	client := binance_connector.NewClient(apiKey, secretKey, baseURL)
	DB := db.ConnectDB(dbURL)
	TGToken := os.Getenv("TGToken")
	Authentication := os.Getenv("Authentication")
	ChatIDStr := os.Getenv("ChatID")
	ChatID, err := strconv.ParseInt(ChatIDStr, 10, 64)
	if err != nil {
		fmt.Println("Failed to parse ChatID:", err)
		return
	}
	bot, err := tg.Connect_tg(TGToken)
	if err != nil {
		fmt.Println("Failed to connect TG:", err)
	}
	go trade.Trade(bot, ChatID, DB, client)
	go om.InitRouter(DB, client, Authentication)
	select {}
}