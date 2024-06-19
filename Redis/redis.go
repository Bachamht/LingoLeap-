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

func main() {
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    initializeUserWordList(rdb)


    // 示例：随机挑选未记住的单词
    words, err := getRandomUnrememberedWords(rdb, 1, 5)
    if err != nil {
        log.Fatal(err)
    }

    for _, word := range words {
        fmt.Println(word)
        markWordAsRemembered(rdb, 1, word)
    }
    

}

 func initializeUserWordList(rdb *redis.Client){
    // 打开词库文件
    file, err := os.Open("./IELTSWords.txt") 
    
    if err != nil {
        log.Fatalf("无法打开文件: %v", err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    userID := 1
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


func markWordAsRemembered(rdb *redis.Client, userID int, word string) error {
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

func getRandomUnrememberedWords(rdb *redis.Client, userID int, numWords int) ([]string, error) {
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