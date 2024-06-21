package api

import (
	"lingoLeap/model"
	"github.com/gin-gonic/gin"
	"net/http"

)

// 提问
func CreateLearning(c *gin.Context) {

	var requestData model.Create_req
		err := c.ShouldBindJSON(&requestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer, err := model.CreateLearning(&requestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	res := model.Create_res{
        Session_id:	0,
		Answer:	answer,	
	}
	c.JSON(http.StatusOK, res)
}


func Communication(c *gin.Context) {

	var requestData model.Communication_req
	err := c.ShouldBindJSON(&requestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer, err := model.CreateCommunication(&requestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	res := model.Communication_res{
		Content: answer,
	}
	c.JSON(http.StatusOK, res)
}