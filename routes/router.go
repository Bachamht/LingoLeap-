package routes

import (
	"lingoLeap/api"
	"github.com/gin-gonic/gin"

)

func InitRouter() {

	r := gin.Default()
	r.Use(corsMiddleware())
	r.POST("createLearning", api.CreateLearning)
	r.POST("communication", api.Communication)
	r.POST("textToSpeech", api.TestToSpeech)
	r.POST("imageCreate", api.ImageCreate)
	r.POST("newRoleDialaogue", api.RoleDialaogueCreate)
	r.Run(":8877")
}
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") 
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") 
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type") 

		// 处理 OPTIONS 预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}