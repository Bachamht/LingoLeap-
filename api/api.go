package api

import (
	"lingoLeap/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"

)

func CreateLearning(c *gin.Context) {

	var requestData model.Create_req
		err := c.ShouldBindJSON(&requestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("c test:", c)
	fmt.Println("requestDataTest:", requestData)

	answer, words, err := model.CreateLearning(&requestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	res := model.Create_res{
        Session_id:	0,
		Answer:	answer,	
		Words: words,
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

func TestToSpeech(c *gin.Context) {
        var requestData model.TextToSpeechRequest
        if err := c.BindJSON(&requestData); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

		audioURL, err := model.TextToSpeech(&requestData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		res := model.TextToSpeechResponse{
			AudioURL: audioURL,
		}
		c.JSON(http.StatusOK, res)
}

func ImageCreate(c *gin.Context) {
		var requestData model.ImageCreate_request
		if err := c.BindJSON(&requestData); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
		
		imageURL, err := model.ImageCrate(&requestData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		res := model.ImageCreate_response{
			AvatarUrl: imageURL,
		}
		c.JSON(http.StatusOK, res)
}

func RoleDialaogueCreate(c *gin.Context) {
	var requestData model.RoleDialaogue_req
		err := c.ShouldBindJSON(&requestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer,  err := model.CreateNewRoleDialaogue(&requestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	res := model.RoleDialaogue_res{
        Session_id:	0,
		Message: answer,	
	}
	c.JSON(http.StatusOK, res)
}