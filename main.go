package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {

	redisCtx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	rdb.ZAdd(redisCtx, "LBSet",
		redis.Z{Score: 0, Member: "8080"},
		redis.Z{Score: 0, Member: "8081"},
		redis.Z{Score: 0, Member: "8082"},
		redis.Z{Score: 0, Member: "8083"},
	)

	r := gin.Default()
	r.GET("/gethandler", func(ctx *gin.Context) {

		res, err := rdb.ZRangeWithScores(ctx, "LBSet", 0, 0).Result()
		if err != nil {
			panic(err)
		}
		port := res[0].Member.(string)
		fmt.Println("Picked server: " + port)
		newScore, err := rdb.ZIncrBy(ctx, "LBSet", 1, port).Result()
		resp, err := http.Get("http://localhost:" + port + "/heartbeat")
		if err != nil {
			panic(err)
		}
		fmt.Println("Port: " + port + " has now score of : " + strconv.Itoa(int(newScore)))
		if err != nil {
			panic(err)
		}

		if err != nil {
			ctx.JSON(resp.StatusCode, gin.H{
				"error": resp.Status,
			})
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		_, decErr := rdb.ZIncrBy(redisCtx, "LBSet", -1, port).Result()
		if decErr != nil {
			panic(decErr)
		}
		ctx.JSON(resp.StatusCode, gin.H{
			"body": result["message"],
		})

	})
	r.Run(":9090")
}
