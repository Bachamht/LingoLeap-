package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/websocket"
)

var (
    db  *sql.DB
    ctx = context.Background()
)

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// 初始化 MySQL 数据库连接
func initDB() {
    var err error
    dsn := "username:password@tcp(localhost:3306)/dbname"
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("无法连接到数据库: %v", err)
    }

    // 验证连接
    if err = db.Ping(); err != nil {
        log.Fatalf("无法连接到数据库: %v", err)
    }
}

func main() {
    initDB()
    defer db.Close()

    d := websocket.Dialer{
        HandshakeTimeout: 5 * time.Second,
    }

    conn, resp, err := d.Dial(assembleAuthUrl1(hostUrl, apiKey, apiSecret), nil)
    if err != nil {
        panic(readResp(resp) + err.Error())
        return
    } else if resp.StatusCode != 101 {
        panic(readResp(resp) + err.Error())
    }

    userID := "example_user_id"
    question := "你是谁，可以干什么？"

    // 获取历史记录
    history, err := getSessionHistory(userID)
    if err != nil {
        log.Fatalf("无法获取历史记录: %v", err)
    }

    go func() {
        data := genParams1(appid, question, history)
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
            return
        }

        payload := data["payload"].(map[string]interface{})
        choices := payload["choices"].(map[string]interface{})
        header := data["header"].(map[string]interface{})
        code := header["code"].(float64)

        if code != 0 {
            fmt.Println(data["payload"])
            return
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

    fmt.Println(answer)

    // 保存会话记录到数据库
    saveSessionRecord(userID, "user", question)
    saveSessionRecord(userID, "assistant", answer)

    time.Sleep(1 * time.Second)
}

// 获取历史会话记录
func getSessionHistory(userID string) ([]Message, error) {
    query := "SELECT role, content FROM session_info WHERE user_id = ? ORDER BY create_timestamp ASC"
    rows, err := db.QueryContext(ctx, query, userID)
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
func saveSessionRecord(userID, role, content string) {
    query := "INSERT INTO session_info (user_id, role, content, create_timestamp, update_timestamp) VALUES (?, ?, ?, NOW(), NOW())"
    _, err := db.Exec(query, userID, role, content)
    if err != nil {
        log.Fatalf("无法保存会话记录: %v", err)
    }
}

// 组装认证 URL
func assembleAuthUrl1(hostUrl, apiKey, apiSecret string) string {
    // 根据实际情况修改返回的数据结构和字段名
    return fmt.Sprintf("%s?apiKey=%s&apiSecret=%s", hostUrl, apiKey, apiSecret)
}

// 读取响应
func readResp(resp *http.Response) string {
    // 根据实际情况修改返回的数据结构和字段名
    return "response"
}

