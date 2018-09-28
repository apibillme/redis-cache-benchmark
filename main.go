package main

import (
	"fmt"
	"time"

	"github.com/spf13/cast"

	"github.com/go-redis/redis"
	vegeta "github.com/tsenart/vegeta/lib"
)

func main() {
	// setup client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	client.Set("123", "555", time.Duration(1*time.Second))

	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 5 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:8000/redis?key=123",
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("Total Requests: %s\n", cast.ToString(metrics.Requests))
	fmt.Printf("Success Ratio: %s\n", cast.ToString(metrics.Success*100)+"%")
	fmt.Printf("Max: %s\n", metrics.Latencies.Max)
	fmt.Printf("Mean: %s\n", metrics.Latencies.Mean)
	fmt.Printf("50th percentile: %s\n", metrics.Latencies.P50)
	fmt.Printf("95th percentile: %s\n", metrics.Latencies.P95)
	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}
