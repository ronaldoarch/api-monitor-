package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"api-monitor/internal/api"
	"api-monitor/internal/monitor"
	"api-monitor/internal/storage"
)

func main() {
	cliMode := flag.Bool("cli", false, "Modo CLI para executar testes")
	portFlag := flag.String("port", "", "Porta para o servidor HTTP")
	flag.Parse()

	if *cliMode {
		runCLI()
		return
	}

	// Railway e outras plataformas cloud usam a variável PORT
	port := os.Getenv("PORT")
	if port == "" {
		if *portFlag != "" {
			port = *portFlag
		} else {
			port = "8080"
		}
	}

	if *cliMode {
		runCLI()
		return
	}

	// Inicializar storage
	store := storage.NewMemoryStorage()

	// Inicializar monitor
	monitorService := monitor.NewMonitor(store)

	// Inicializar API
	apiHandler := api.NewAPI(monitorService, store)

	// Configurar rotas
	router := apiHandler.SetupRoutes()
	
	// Servir arquivos estáticos
	fs := http.FileServer(http.Dir("./web/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Dashboard
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	// WebSocket para atualizações em tempo real
	router.HandleFunc("/ws", apiHandler.HandleWebSocket)
	
	http.Handle("/", router)

	fmt.Printf("🚀 API Monitor rodando em http://localhost:%s\n", port)
	fmt.Printf("📊 Dashboard: http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func runCLI() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: go run main.go --cli <comando>")
		fmt.Println("Comandos disponíveis:")
		fmt.Println("  test <url> - Executa um teste rápido")
		fmt.Println("  load <url> <requests> <concurrency> - Executa teste de carga")
		os.Exit(1)
	}

	command := os.Args[2]
	store := storage.NewMemoryStorage()
	monitorService := monitor.NewMonitor(store)

	switch command {
	case "test":
		if len(os.Args) < 4 {
			fmt.Println("Uso: go run main.go --cli test <url>")
			os.Exit(1)
		}
		url := os.Args[3]
		fmt.Printf("Executando teste rápido em: %s\n", url)
		result := monitorService.RunQuickTest(url)
		fmt.Printf("\n✅ Resultado do Teste:\n")
		fmt.Printf("  Status: %d\n", result.Status)
		fmt.Printf("  Sucesso: %v\n", result.Success)
		fmt.Printf("  Tempo de Resposta: %dms\n", result.Duration)
		fmt.Printf("  Tamanho da Resposta: %d bytes\n", result.ResponseSize)
		if result.Error != "" {
			fmt.Printf("  Erro: %s\n", result.Error)
		}
	case "load":
		if len(os.Args) < 6 {
			fmt.Println("Uso: go run main.go --cli load <url> <requests> <concurrency>")
			os.Exit(1)
		}
		url := os.Args[3]
		var requests, concurrency int
		fmt.Sscanf(os.Args[4], "%d", &requests)
		fmt.Sscanf(os.Args[5], "%d", &concurrency)
		fmt.Printf("Executando teste de carga:\n")
		fmt.Printf("  URL: %s\n", url)
		fmt.Printf("  Requisições: %d\n", requests)
		fmt.Printf("  Concorrência: %d\n", concurrency)
		fmt.Println("\n⏳ Executando...")
		result := monitorService.RunLoadTest(url, requests, concurrency)
		fmt.Printf("\n✅ Resultado do Teste de Carga:\n")
		fmt.Printf("  Total de Requisições: %d\n", result.TotalRequests)
		fmt.Printf("  Sucesso: %d\n", result.SuccessCount)
		fmt.Printf("  Erros: %d\n", result.ErrorCount)
		fmt.Printf("  Taxa de Sucesso: %.2f%%\n", float64(result.SuccessCount)/float64(result.TotalRequests)*100)
		fmt.Printf("  Tempo Médio de Resposta: %.2fms\n", result.AvgResponseTime)
		fmt.Printf("  Tempo Mínimo: %dms\n", result.MinResponseTime)
		fmt.Printf("  Tempo Máximo: %dms\n", result.MaxResponseTime)
		fmt.Printf("  Duração Total: %dms\n", result.Duration)
		fmt.Printf("  Códigos de Status:\n")
		for code, count := range result.StatusCodes {
			fmt.Printf("    %d: %d\n", code, count)
		}
	default:
		fmt.Printf("Comando desconhecido: %s\n", command)
		os.Exit(1)
	}
}

