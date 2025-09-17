package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const N = 10

var (
	addr           = flag.String("http-address", "127.0.0.3:8080", "HTTP Host and Port")
	iterations     = flag.Int("iterations", 100, "Number of Iterations for writing")
	readIterations = flag.Int("readIterations", 100, "Number of Iterations for reading")
	concurrency    = flag.Int("concurrency", 5, "Number of concurrnt go functions")
)
var httpClient = &http.Client{
	Transport: &http.Transport{
		IdleConnTimeout:     time.Second * 60,
		MaxIdleConns:        32,
		MaxConnsPerHost:     32,
		MaxIdleConnsPerHost: 32,
	}}

func benchmark(funcname string, iterations int, fn func() string) (qps float64, strs []string) { //can calculate min and max also
	var max time.Duration
	var min time.Duration

	start := time.Now()
	for i := 0; i < iterations; i++ {
		strs = append(strs, fn())
		iterTime := time.Since(start)
		if iterTime > max {
			max = iterTime
		}
		if iterTime < min {
			max = iterTime
		}
	}
	avg := time.Since(start) / N
	qps = float64(iterations) / (float64(time.Since(start)) / float64(time.Second))
	fmt.Printf("func %s- avg: %s, min: %s, max: %s, QPS: %.1f\n", funcname, avg, min, max, qps)
	return qps, strs
}

func benchmarkWrite() (strs []string) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalQPS float64
	var allKeys []string
	for i := 0; i < *concurrency; i += 1 {
		wg.Add(1)
		go func() {
			qps, strs := benchmark("write", *iterations, randTestWrite)
			mu.Lock()
			totalQPS += qps
			allKeys = append(allKeys, strs...)
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	log.Printf("Total QPS= %.1f, set %d keys", totalQPS, len(allKeys))
	return allKeys
}

func randTestWrite() (key string) {
	key = fmt.Sprintf("key-%d", rand.Intn(10000)) //no need to check quality of key
	value := fmt.Sprintf("value-%d", rand.Intn(10000))
	values := url.Values{}
	values.Set("key", key)
	values.Set("value", value)
	resp, err := httpClient.Get("http://" + *addr + "/set?" + values.Encode())
	if err != nil {
		log.Fatal("Could not GET")
	}
	io.Copy(io.Discard, resp.Body)
	defer resp.Body.Close()
	return key
}
func benchmarkRead(allKeys []string) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalQPS float64
	for i := 0; i < *concurrency; i += 1 {
		wg.Add(1)
		go func() {
			qps, _ := benchmark("read", *readIterations, func() string { return randTestRead(allKeys) })
			mu.Lock()
			totalQPS += qps
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	log.Printf("Total QPS= %.1f, set %d keys", totalQPS, len(allKeys))
}

func randTestRead(allKeys []string) (key string) {
	key = allKeys[rand.Intn(len(allKeys))]
	values := url.Values{}
	values.Set("key", key)
	resp, err := httpClient.Get("http://" + *addr + "/get?" + values.Encode())
	if err != nil {
		log.Fatal("Could not GET")
	}
	io.Copy(io.Discard, resp.Body)
	defer resp.Body.Close()
	return key
}

func main() {
	flag.Parse()
	allKeys := benchmarkWrite()
	benchmarkRead(allKeys)
}
