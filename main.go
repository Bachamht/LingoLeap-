package main

import (
	"lingoLeap/model"
	"lingoLeap/Redis"
)

func main() {
	model.ConnectMysql()
	Redis.ConnectRedis()
	model.ConnetcSpark()
}