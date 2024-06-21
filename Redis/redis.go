package main

import (
    "context"
    "fmt"
    "github.com/go-redis/redis/v8"
    "log"
    "math/rand"
    "time"
    "regexp"
    "bufio"
    "os"
    
)

var ctx = context.Background()
var rbd *redis.Client

func ConnectRedis() {
    rdb = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
}


func CreateLearning(userID string, numWords int) ([]string, error) {
    // 检查 Redis 中是否已经存在该用户的单词库
    pattern := fmt.Sprintf("word:%s:*", userID)
    keys, err := rdb.Keys(ctx, pattern).Result()
    if err != nil {
        return nil, err
    }

    if len(keys) == 0 {
        err = initializeUserWordList(rdb, userID)
        if err != nil {
            return nil, err
        }
    }

    // 获取指定数量的未记住单词
    words, err := getRandomUnrememberedWords(rdb, userID, numWords)
    if err != nil {
        return nil, err
    }

    return words, nil
}

 func initializeUserWordList(userID int){
    // 打开词库文件
    file, err := os.Open("./IELTSWords.txt") 
    
    if err != nil {
        log.Fatalf("无法打开文件: %v", err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    wordID := 1

    re := regexp.MustCompile(`^[a-zA-Z]+`)

    for scanner.Scan() {
        line := scanner.Text()
        matches := re.FindStringSubmatch(line)
        if len(matches) > 0 {
            word := matches[0]
            key := fmt.Sprintf("word:%d:%d", userID, wordID)
            err := rdb.HMSet(ctx, key, map[string]interface{}{
                "content":      word,
                "isRemembered": 0,
            }).Err()
            if err != nil {
                log.Fatalf("无法存储到 Redis: %v", err)
            }
            wordID++
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatalf("扫描文件时出错: %v", err)
    }

    fmt.Println("单词已成功导入 Redis")
}


func markWordAsRemembered(userID int, word string) error {
    pattern := fmt.Sprintf("word:%d:*", userID)
    keys, err := rdb.Keys(ctx, pattern).Result()
    if err != nil {
        return err
    }

    for _, key := range keys {
        content, err := rdb.HGet(ctx, key, "content").Result()
        if err != nil {
            return err
        }
        if content == word {
            fmt.Printf("标记单词 %s\n", word, key)
            err = rdb.HSet(ctx, key, "isRemembered", "1").Err()
            if err != nil {
                return err
            }
            break
        }
    }

    return nil
}

func getRandomUnrememberedWords(userID int, numWords int) ([]string, error) {
    pattern := fmt.Sprintf("word:%d:*", userID)
    keys, err := rdb.Keys(ctx, pattern).Result()
    if err != nil {
        return nil, err
    }

    var unrememberedKeys []string
    for _, key := range keys {
        isRemembered, err := rdb.HGet(ctx, key, "isRemembered").Result()
        if err != nil {
            return nil, err
        }
        if isRemembered == "0" {
            unrememberedKeys = append(unrememberedKeys, key)
        }
    }

    rand.Seed(time.Now().UnixNano())
    selectedKeys := make(map[int]struct{})
    var words []string
    for len(words) < numWords && len(selectedKeys) < len(unrememberedKeys) {
        idx := rand.Intn(len(unrememberedKeys))
        if _, exists := selectedKeys[idx]; !exists {
            selectedKeys[idx] = struct{}{}
            content, err := rdb.HGet(ctx, unrememberedKeys[idx], "content").Result()
            if err != nil {
                return nil, err
            }
            words = append(words, content)
        }
    }

    return words, nil
}