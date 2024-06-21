
package model
import (
	"fmt"
)

type Communication_req struct {
	User_id     int `json:"user_id"`
	Session_id int    `json:"session_id"`
	Content string `json:"content"`
}

type Communication_res struct {
	Content string `json:"content"`
}


func CreateCommunication(requestData *Communication_req) (string, error) {
	userID := requestData.User_id
	//sessionID := requestData.Session_id
	content := requestData.Content
	fmt.Println(content)
	role := "user"
	answer, err := IteracionWithAI(content, userID, role)
	if err != nil {
		return "", err
	}
	return answer, err
}