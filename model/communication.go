
package model
import (
	"fmt"
	"log"
	"time"
	"encoding/json"
	"net/http"
	"github.com/gorilla/websocket"
)

type Communication_req struct {
	User_id     string `json:"user_id"`
	Session_id int    `json:"session_id"`
	Content string `json:"content"`
}

type Communication_res struct {
	Content string `json:"content"`
}


func CreateLearning(requestData *Communication_req) (string, error) {
	userID := requestData.User_id
	sessionID := requestData.Session_id
	content := requestData.Content
	role := "user"
	asnwer, err := IteracionWithAI(content, userID, role)
	if err != nil {
		return "", err
	}
	return answer, err
}