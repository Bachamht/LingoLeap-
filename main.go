package main

import (
	"lingoLeap/model"
	"lingoLeap/redis"
	"lingoLeap/routes"
	"fmt"
)

func main() {
	model.ConnectMysql()
	fmt.Println("Connected to database successfully.")
	redis.ConnectRedis()
	fmt.Println("Connected to Redis successfully.")
	model.ConnetcSpark()
	fmt.Println("Connected to Spark successfully.")
	routes.InitRouter()
	fmt.Println("Server started successfully.")
}