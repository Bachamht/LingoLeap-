
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
    Words      []string
}

func CreateLearning(requestData *Create_req) (string, []string, error){
	userID := requestData.User_id
    wordNumber := requestData.Word_number
	words, err := redis.CreateLearning(userID, wordNumber)
    fmt.Println("userID", userID)
    fmt.Println("wordNumber", wordNumber)

	prompt := "你的角色是一个精通中英文的英语老师，用户会给你发送一些生词，你帮助用户以对话的形式记忆单词。\n\n" +
        "Step1:你要给出生词的常见中文注释，可以适当拓展（一个单词可能有多个注释）\n\n" +
        "Step2:你要用这几个单词构建一个或者多个简短的英文小故事，并给出中文翻译\n\n" +
        "举例：你收到三个单词：apple environment hacker\n" +
        "回复：apple  n.苹果\n\n" +
        "environment  n.自然环境，生态环境；周围状况，条件；工作平台，软件包\n\n" +
        "hacker   n.黑客，骇客；不擅长某项运动的人；计算机迷；砍（或劈）的人，用于砍（或劈）的东西\n\n" +
        "Tom loves apples. He lives in a town with a beautiful **environment**. One day, while eating an **apple**, he sees a **hacker** near the apple trees. The **hacker** is using a computer to help the apple trees grow better and protect the **environment**.\n\n" +
        "汤姆喜欢苹果。他住在一个环境优美的小镇上。一天，当他在吃苹果的时候，他看到苹果树旁边有一个黑客。**黑客**正在用电脑帮助苹果树长得更好，保护环境。\n\n" +
        "本次你收到的单词为：\n\n"
    for _, word := range words {
        prompt += fmt.Sprintf("- %s\n", word)
    }
	role := "system"

	answer, err1 := IteracionWithAI(prompt, userID, role)
    if err1 != nil {
        return "", nil, err
    }

	return answer, words, nil
}