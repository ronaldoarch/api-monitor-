package storage

import (
	"sync"
	"time"
)

type TestResult struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Method      string    `json:"method"`
	Status      int       `json:"status"`
	Duration    int64     `json:"duration"` // em milissegundos
	Timestamp   time.Time `json:"timestamp"`
	Success     bool      `json:"success"`
	Error       string    `json:"error,omitempty"`
	ResponseSize int64    `json:"response_size"`
}

type LoadTestResult struct {
	ID              string       `json:"id"`
	URL             string       `json:"url"`
	TotalRequests   int          `json:"total_requests"`
	Concurrency     int          `json:"concurrency"`
	Duration        int64        `json:"duration"` // em milissegundos
	SuccessCount    int          `json:"success_count"`
	ErrorCount      int          `json:"error_count"`
	AvgResponseTime float64      `json:"avg_response_time"`
	MinResponseTime int64        `json:"min_response_time"`
	MaxResponseTime int64        `json:"max_response_time"`
	StatusCodes     map[int]int  `json:"status_codes"`
	Timestamp       time.Time    `json:"timestamp"`
	Results         []TestResult `json:"results"`
}

type Storage interface {
	SaveTestResult(result TestResult)
	SaveLoadTestResult(result LoadTestResult)
	GetTestResults(limit int) []TestResult
	GetLoadTestResults(limit int) []LoadTestResult
	GetLoadTestByID(id string) (*LoadTestResult, bool)
}

type MemoryStorage struct {
	mu              sync.RWMutex
	testResults     []TestResult
	loadTestResults []LoadTestResult
	loadTestsByID   map[string]*LoadTestResult
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		testResults:     make([]TestResult, 0),
		loadTestResults: make([]LoadTestResult, 0),
		loadTestsByID:   make(map[string]*LoadTestResult),
	}
}

func (s *MemoryStorage) SaveTestResult(result TestResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.testResults = append(s.testResults, result)
	if len(s.testResults) > 1000 {
		s.testResults = s.testResults[1:]
	}
}

func (s *MemoryStorage) SaveLoadTestResult(result LoadTestResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.loadTestResults = append(s.loadTestResults, result)
	s.loadTestsByID[result.ID] = &result
	if len(s.loadTestResults) > 100 {
		s.loadTestResults = s.loadTestResults[1:]
	}
}

func (s *MemoryStorage) GetTestResults(limit int) []TestResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if limit > len(s.testResults) {
		limit = len(s.testResults)
	}
	start := len(s.testResults) - limit
	if start < 0 {
		start = 0
	}
	result := make([]TestResult, limit)
	copy(result, s.testResults[start:])
	return result
}

func (s *MemoryStorage) GetLoadTestResults(limit int) []LoadTestResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if limit > len(s.loadTestResults) {
		limit = len(s.loadTestResults)
	}
	start := len(s.loadTestResults) - limit
	if start < 0 {
		start = 0
	}
	result := make([]LoadTestResult, limit)
	copy(result, s.loadTestResults[start:])
	return result
}

func (s *MemoryStorage) GetLoadTestByID(id string) (*LoadTestResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, ok := s.loadTestsByID[id]
	return result, ok
}

