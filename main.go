package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/plivo/decimal"
	progressbar "github.com/schollz/progressbar/v3"
)

// it's totally possible to strip this down further to get at the
// bare minumum reproducer but this is good enough for our use

var (
	redisAddr       = "127.0.0.1:30001"
	concurrentIters = 999999
	redisMaxRetries = 2
	redisClientName = "gotham"
	cacheExpiry     = time.Duration(3600) * time.Second
	dcr             = decimal.NewFromFloat(42.24)
	dtr             = decimal.NewFromFloat(24.42)

	client *redis.ClusterClient
)

func main() {
	client = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:      []string{redisAddr},
		MaxRetries: redisMaxRetries,
		PoolSize:   20,
		OnConnect: func(conn *redis.Conn) error {
			return conn.ClientSetName(redisClientName).Err()
		},
		DialTimeout:  10 * time.Millisecond,
		ReadTimeout:  10 * time.Millisecond,
		WriteTimeout: 10 * time.Millisecond,
	})

	bar := progressbar.Default(int64(concurrentIters))
	wg := &sync.WaitGroup{}
	for i := 0; i < concurrentIters; i++ {
		wg.Add(2)
		bar.Add(1)
		go setCache(wg)
		go getCache(wg)
	}

	wg.Wait()
}

func setCache(wg *sync.WaitGroup) error {
	defer wg.Done()

	rate := make(map[string]decimal.Decimal)
	rate["carrier_rate"] = dcr
	rate["termination_rate"] = dtr

	rateStr, err := json.Marshal(rate)
	if err != nil {
		return fmt.Errorf("json.Marshal() failed")
	}

	_, err = client.Set("key", string(rateStr), cacheExpiry).Result()
	return err
}

func getCache(wg *sync.WaitGroup) error {
	defer wg.Done()

	rate, err := client.Get("key").Result()
	if err != nil {
		return err
	}

	if len(rate) == 0 {
		return fmt.Errorf("some error")
	}

	type smsCredit struct {
		CarrierAmount   decimal.Decimal `json:"carrier_amount"`
		CarrierRate     decimal.Decimal `json:"carrier_rate"`
		TerminationRate decimal.Decimal `json:"termination_rate"`
		SmsGroupId      int             `json:"sms_group_id"`
	}

	var response smsCredit
	err = json.Unmarshal([]byte(rate), &response)
	if err != nil {
		log.Fatalf("\njson.Unmarshal() failed: %s\n", err.Error()) // we hit the bug
		return err
	}

	return nil
}
