package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

type TextToSpeechRequest struct {
    Text string `json:"text"`
}

type TextToSpeechResponse struct {
    AudioURL string `json:"audioUrl"`
}

type WsParam struct {
    APPID     string
    APIKey    string
    APISecret string
    Text      string
}

func (w *WsParam) CreateURL() (string, error) {
    baseURL := "wss://tts-api.xfyun.cn/v2/tts"
    now := time.Now()
    date := now.Format(time.RFC1123)

    signatureOrigin := fmt.Sprintf("host: ws-api.xfyun.cn\ndate: %s\nGET /v2/tts HTTP/1.1", date)
    mac := hmac.New(sha256.New, []byte(w.APISecret))
    mac.Write([]byte(signatureOrigin))
    signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

    authorizationOrigin := fmt.Sprintf(
        `api_key="%s", algorithm="hmac-sha256", headers="host date request-line", signature="%s"`,
        w.APIKey, signature)
    authorization := base64.StdEncoding.EncodeToString([]byte(authorizationOrigin))

    v := url.Values{}
    v.Add("authorization", authorization)
    v.Add("date", date)
    v.Add("host", "ws-api.xfyun.cn")

    return baseURL + "?" + v.Encode(), nil
}

func websocketConnectAndReceive(wsURL string) ([]byte, error) {
    c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err != nil {
        return nil, fmt.Errorf("dial error: %v", err)
    }
    defer c.Close()

    var audioData []byte
    done := make(chan struct{})

    go func() {
        defer close(done)
        for {
            _, message, err := c.ReadMessage()
            if err != nil {
                fmt.Println("read error:", err)
                return
            }

            var resp map[string]interface{}
            err = json.Unmarshal(message, &resp)
            if err != nil {
                fmt.Println("json unmarshal error:", err)
                continue
            }

            if resp["code"].(float64) != 0 {
                fmt.Printf("error code: %v, message: %v\n", resp["code"], resp["message"])
                continue
            }

            data := resp["data"].(map[string]interface{})
            audio, err := base64.StdEncoding.DecodeString(data["audio"].(string))
            if err != nil {
                fmt.Println("base64 decode error:", err)
                continue
            }

            audioData = append(audioData, audio...)
            if data["status"].(float64) == 2 {
                break
            }
        }
    }()

    <-done
    return audioData, nil
}

func main() {
    router := gin.Default()

    router.POST("/texttospeech", func(c *gin.Context) {
        var req TextToSpeechRequest
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        wsParam := WsParam{
            APPID:     "your_app_id",
            APIKey:    "your_api_key",
            APISecret: "your_api_secret",
            Text:      req.Text,
        }

        wsURL, err := wsParam.CreateURL()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        audioData, err := websocketConnectAndReceive(wsURL)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        audioFilePath := "./static/audio/demo.pcm"
        err = ioutil.WriteFile(audioFilePath, audioData, 0644)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        audioURL := "http://yourdomain.com/static/audio/demo.pcm"
        c.JSON(http.StatusOK, TextToSpeechResponse{AudioURL: audioURL})
    })

    router.Static("/static", "./static")

    router.Run
