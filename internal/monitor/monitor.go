package monitor

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"api-monitor/internal/storage"

	"github.com/google/uuid"
)

type Monitor struct {
	storage storage.Storage
	client  *http.Client
}

func NewMonitor(storage storage.Storage) *Monitor {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	return &Monitor{
		storage: storage,
		client:  client,
	}
}

func (m *Monitor) RunQuickTest(url string) *storage.TestResult {
	start := time.Now()
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &storage.TestResult{
			ID:        uuid.New().String(),
			URL:       url,
			Method:    "GET",
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now(),
		}
	}

	resp, err := m.client.Do(req)
	duration := time.Since(start).Milliseconds()

	result := storage.TestResult{
		ID:          uuid.New().String(),
		URL:         url,
		Method:      "GET",
		Duration:    duration,
		Timestamp:   time.Now(),
		Success:     err == nil && resp != nil && resp.StatusCode < 400,
		ResponseSize: 0,
	}

	if err != nil {
		result.Error = err.Error()
		result.Status = 0
	} else {
		result.Status = resp.StatusCode
		if resp.Body != nil {
			body, _ := io.ReadAll(resp.Body)
			result.ResponseSize = int64(len(body))
			resp.Body.Close()
		}
	}

	m.storage.SaveTestResult(result)
	return &result
}

func (m *Monitor) RunLoadTest(url string, totalRequests, concurrency int) *storage.LoadTestResult {
	start := time.Now()
	testID := uuid.New().String()

	results := make([]storage.TestResult, 0, totalRequests)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Canal para controlar concorrência
	semaphore := make(chan struct{}, concurrency)
	requestChan := make(chan int, totalRequests)

	// Preencher canal de requisições
	for i := 0; i < totalRequests; i++ {
		requestChan <- i
	}
	close(requestChan)

	// Executar requisições concorrentes
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range requestChan {
				semaphore <- struct{}{}
				result := m.RunQuickTest(url)
				result.ID = fmt.Sprintf("%s-%d", testID, len(results))
				mu.Lock()
				results = append(results, *result)
				mu.Unlock()
				<-semaphore
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start).Milliseconds()

	// Calcular estatísticas
	successCount := 0
	errorCount := 0
	var totalDuration int64
	var minDuration int64 = 999999999
	var maxDuration int64
	statusCodes := make(map[int]int)

	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			errorCount++
		}
		totalDuration += result.Duration
		if result.Duration < minDuration {
			minDuration = result.Duration
		}
		if result.Duration > maxDuration {
			maxDuration = result.Duration
		}
		statusCodes[result.Status]++
	}

	avgResponseTime := float64(totalDuration) / float64(len(results))

	loadTestResult := storage.LoadTestResult{
		ID:              testID,
		URL:             url,
		TotalRequests:   totalRequests,
		Concurrency:     concurrency,
		Duration:        duration,
		SuccessCount:    successCount,
		ErrorCount:      errorCount,
		AvgResponseTime: avgResponseTime,
		MinResponseTime: minDuration,
		MaxResponseTime: maxDuration,
		StatusCodes:     statusCodes,
		Timestamp:       time.Now(),
		Results:         results,
	}

	m.storage.SaveLoadTestResult(loadTestResult)
	return &loadTestResult
}

