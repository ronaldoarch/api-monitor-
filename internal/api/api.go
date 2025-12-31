package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"api-monitor/internal/monitor"
	"api-monitor/internal/storage"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type API struct {
	monitor *monitor.Monitor
	storage storage.Storage
	upgrader websocket.Upgrader
	clients map[*websocket.Conn]bool
	broadcast chan []byte
}

func NewAPI(monitor *monitor.Monitor, storage storage.Storage) *API {
	api := &API{
		monitor: monitor,
		storage: storage,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients: make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
	go api.handleBroadcast()
	return api
}

func (a *API) SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()
	
	apiRouter.HandleFunc("/test", a.handleTest).Methods("POST")
	apiRouter.HandleFunc("/load", a.handleLoadTest).Methods("POST")
	apiRouter.HandleFunc("/results", a.handleGetResults).Methods("GET")
	apiRouter.HandleFunc("/load-results", a.handleGetLoadResults).Methods("GET")
	apiRouter.HandleFunc("/load-results/{id}", a.handleGetLoadResultByID).Methods("GET")
	
	return router
}

func (a *API) handleTest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	result := a.monitor.RunQuickTest(req.URL)
	a.broadcastUpdate("test_result", result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (a *API) handleLoadTest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL         string `json:"url"`
		Requests    int    `json:"requests"`
		Concurrency int    `json:"concurrency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	if req.Requests <= 0 {
		req.Requests = 100
	}
	if req.Concurrency <= 0 {
		req.Concurrency = 10
	}

	// Executar em goroutine para não bloquear
	go func() {
		result := a.monitor.RunLoadTest(req.URL, req.Requests, req.Concurrency)
		a.broadcastUpdate("load_test_result", result)
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "started",
		"message": "Load test iniciado",
	})
}

func (a *API) handleGetResults(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	results := a.storage.GetTestResults(limit)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (a *API) handleGetLoadResults(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	results := a.storage.GetLoadTestResults(limit)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (a *API) handleGetLoadResultByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, ok := a.storage.GetLoadTestByID(id)
	if !ok {
		http.Error(w, "Test not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (a *API) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	a.clients[conn] = true

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(a.clients, conn)
			break
		}
	}
}

func (a *API) broadcastUpdate(eventType string, data interface{}) {
	message := map[string]interface{}{
		"type": eventType,
		"data": data,
		"timestamp": time.Now(),
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return
	}

	select {
	case a.broadcast <- jsonData:
	default:
	}
}

func (a *API) handleBroadcast() {
	for {
		message := <-a.broadcast
		for client := range a.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				delete(a.clients, client)
				client.Close()
			}
		}
	}
}

