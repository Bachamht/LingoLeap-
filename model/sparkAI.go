package model

import(
	"fmt"
	"log"
	"time"
	"encoding/json"
	"net/http"
	"github.com/gorilla/websocket"
	"os"
	"strings"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"net/url"
	"context"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var (
	conn *websocket.Conn
	resp *http.Response
	err error
	hostUrl string
	appid string
	apiSecret string
	apiKey string
	ctx = context.Background()
)


func ConnetcSpark() {
    d := websocket.Dialer{
        HandshakeTimeout: 5 * time.Second,
    }
	hostUrl = os.Getenv("HostURL")
	apiKey = os.Getenv("APIKEY")
	apiSecret = os.Getenv("APISecret")
	appid = os.Getenv("APPID")

    conn, resp, err = d.Dial(assembleAuthUrl1(hostUrl, apiKey, apiSecret), nil)
    if err != nil {
        panic(readResp(resp) + err.Error())
        return
    } else if resp.StatusCode != 101 {
        panic(readResp(resp) + err.Error())
    }
}

func CheckAndReconnect() {
    if conn == nil || conn.UnderlyingConn().RemoteAddr() == nil {
        fmt.Println("连接无效，重新连接...")
        ConnectSpark()
    }
}

func IteracionWithAI(prompt string, userID int, role string) (string, error){
	CheckAndReconnect()
    history, err := getSessionHistory(userID)
    if err != nil {
        log.Fatalf("无法获取历史记录: %v", err)
    }

    go func() {
        data := genParams1(appid, prompt, history)
        conn.WriteJSON(data)
    }()

    var answer = ""
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("读取消息错误:", err)
            break
        }

        var data map[string]interface{}
        err1 := json.Unmarshal(msg, &data)
        if err1 != nil {
            fmt.Println("解析 JSON 错误:", err)
            return "", err1
        }

        payload := data["payload"].(map[string]interface{})
        choices := payload["choices"].(map[string]interface{})
        header := data["header"].(map[string]interface{})
        code := header["code"].(float64)

        if code != 0 {
            fmt.Println(data["payload"])
            return "", nil
        }
        status := choices["status"].(float64)
        text := choices["text"].([]interface{})
        content := text[0].(map[string]interface{})["content"].(string)
        if status != 2 {
            answer += content
        } else {
            fmt.Println("收到最终结果")
            answer += content
            usage := payload["usage"].(map[string]interface{})
            temp := usage["text"].(map[string]interface{})
            totalTokens := temp["total_tokens"].(float64)
            fmt.Println("total_tokens:", totalTokens)
            conn.Close()
            break 
        }
    }

    // 保存会话记录到数据库
    saveSessionRecord(userID, role, prompt)
    saveSessionRecord(userID, "assistant", answer)

	fmt.Println(answer)
	return answer, nil
}

// 获取历史会话记录
func getSessionHistory(userID int) ([]Message, error) {
    query := "SELECT role, content FROM session_info WHERE user_id = ? ORDER BY create_timestamp ASC"
    rows, err := DB.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var history []Message
    for rows.Next() {
        var role, content string
        if err := rows.Scan(&role, &content); err != nil {
            return nil, err
        }
        history = append(history, Message{Role: role, Content: content})
    }

    return history, nil
}

// 生成参数
func genParams1(appid, question string, history []Message) map[string]interface{} {
    // 在历史记录中添加最新的问题
    history = append(history, Message{Role: "user", Content: question})

    data := map[string]interface{}{
        "header": map[string]interface{}{
            "app_id": appid,
        },
        "parameter": map[string]interface{}{
            "chat": map[string]interface{}{
                "domain":      "generalv3.5",
                "temperature": float64(0.5),
                "max_tokens":  int64(1024),
            },
        },
        "payload": map[string]interface{}{
            "message": map[string]interface{}{
                "text": history,
            },
        },
    }
    return data
}

// 保存会话记录到数据库
func saveSessionRecord(userID int, role, content string) {
    query := "INSERT INTO session_info (user_id, role, content, create_timestamp, update_timestamp) VALUES (?, ?, ?, NOW(), NOW())"
    _, err := DB.Exec(query, userID, role, content)
    if err != nil {
        log.Fatalf("无法保存会话记录: %v", err)
    }
}



// 创建鉴权url  apikey 即 hmac username
func assembleAuthUrl1(hosturl string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hosturl)
	if err != nil {
		fmt.Println(err)
	}
	//签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	//date = "Tue, 28 May 2019 09:10:42 MST"
	//参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	//拼接签名字符串
	sgin := strings.Join(signString, "\n")
	// fmt.Println(sgin)
	//签名结果
	sha := HmacWithShaTobase64("hmac-sha256", sgin, apiSecret)
	// fmt.Println(sha)
	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callurl := hosturl + "?" + v.Encode()
	return callurl
}


func HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}

