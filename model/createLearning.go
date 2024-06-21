
package model
import (
	"fmt"
	"lingoLeap/redis"
)

type Create_req struct {
	User_id     int `json:"user_id"`
	Word_number int    `json:"word_number"`
}

type Create_res struct {
	Session_id  int
	Answer	 string
}

func CreateLearning(requestData *Create_req) (string, error){
	userID := requestData.User_id
    wordNumber := requestData.Word_number
	words, err := redis.CreateLearning(userID, wordNumber)

	prompt := "你现在的身份是一个英语老师，以对话的形式帮助学生记忆单词。本次需要记忆的单词有：\n\n"
    for _, word := range words {
        prompt += fmt.Sprintf("- %s\n", word)
    }
    prompt += "\n请你根据这些单词生成一个场景与学生进行英文对话，在对话之前先向学生解释这些单词的意思。"
	role := "system"

	answer, err1 := IteracionWithAI(prompt, userID, role)
    if err1 != nil {
        return "", err
    }

	return answer, nil
}