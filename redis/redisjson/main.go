package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-by-example/redis/client"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
)

// Name - student name
type Name struct {
	First  string `json:"first,omitempty"`
	Middle string `json:"middle,omitempty"`
	Last   string `json:"last,omitempty"`
}

// Student - student object
type Student struct {
	Name Name `json:"name,omitempty"`
	Rank int  `json:"rank,omitempty"`
}

func main() {
	rh := rejson.NewReJSONHandler()
	flag.Parse()

	// GoRedis Client
	cli := client.Client()
	defer func() {
		// if err := cli.FlushAll(context.Background()).Err(); err != nil {
		// 	log.Fatalf("goredis - failed to flush: %v", err)
		// }
		if err := cli.Close(); err != nil {
			log.Fatalf("goredis - failed to communicate to redis-server: %v", err)
		}
	}()
	rh.SetGoRedisClient(cli)
	fmt.Println("\nExecuting Example_JSONSET for GoRedis Client")

	student := Student{
		Name: Name{
			"Mark",
			"S",
			"Pronto",
		},
		Rank: 1,
	}
	res, err := rh.JSONSet("student", ".", student)
	if err != nil {
		log.Fatalf("Failed to JSONSet")
		return
	}

	if res.(string) == "OK" {
		fmt.Printf("Success: %s\n", res)
	} else {
		fmt.Println("Failed to Set: ")
	}

	studentJSON, err := redis.Bytes(rh.JSONGet("student", "."))
	if err != nil {
		log.Fatalf("Failed to JSONGet")
		return
	}

	readStudent := Student{}
	err = json.Unmarshal(studentJSON, &readStudent)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal")
		return
	}

	fmt.Printf("Student read from redis : %#v\n", readStudent)
}
