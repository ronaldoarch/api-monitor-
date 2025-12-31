# 🚀 API Monitor

Uma ferramenta moderna e poderosa de monitoramento e teste de performance de APIs desenvolvida em Go. Perfeita para desenvolvedores que precisam testar, monitorar e analisar o desempenho de suas APIs.

## ✨ Características

- **Teste Rápido**: Execute testes individuais em APIs e obtenha métricas detalhadas
- **Teste de Carga**: Realize testes de carga com controle de concorrência e número de requisições
- **Dashboard Web Moderno**: Interface bonita e intuitiva com atualizações em tempo real
- **WebSocket**: Atualizações em tempo real dos resultados dos testes
- **CLI Tool**: Execute testes via linha de comando para automação
- **Métricas Detalhadas**: Tempo de resposta, taxa de sucesso, códigos de status HTTP e muito mais
- **Histórico**: Visualize todos os testes executados anteriormente

## 🛠️ Tecnologias Utilizadas

- **Go 1.21+**: Linguagem principal
- **Gorilla Mux**: Roteamento HTTP
- **Gorilla WebSocket**: Comunicação em tempo real
- **HTML/CSS/JavaScript**: Interface web moderna
- **Goroutines**: Processamento concorrente para testes de carga

## 📋 Pré-requisitos

- Go 1.21 ou superior
- Navegador web moderno

## 🚀 Instalação

1. Clone o repositório:
```bash
git clone <seu-repositorio>
cd api-monitor
```

2. Instale as dependências:
```bash
go mod download
```

3. Execute o servidor:
```bash
go run main.go
```

Ou especifique uma porta customizada:
```bash
go run main.go -port 3000
```

4. Acesse o dashboard no navegador:
```
http://localhost:8080
```

## 💻 Uso

### Modo Web (Dashboard)

1. Inicie o servidor com `go run main.go`
2. Acesse `http://localhost:8080` no navegador
3. Use as abas para:
   - **Teste Rápido**: Teste uma API individualmente
   - **Teste de Carga**: Execute testes de carga com múltiplas requisições
   - **Histórico**: Visualize testes anteriores

### Modo CLI

Execute testes via linha de comando:

```bash
# Teste rápido
go run main.go --cli test https://api.exemplo.com/endpoint

# Teste de carga
go run main.go --cli load https://api.exemplo.com/endpoint 100 10
# Parâmetros: URL, número de requisições, concorrência
```

## 📡 API REST

### Endpoints Disponíveis

#### POST `/api/test`
Executa um teste rápido em uma API.

**Request Body:**
```json
{
  "url": "https://api.exemplo.com/endpoint"
}
```

**Response:**
```json
{
  "id": "uuid",
  "url": "https://api.exemplo.com/endpoint",
  "method": "GET",
  "status": 200,
  "duration": 150,
  "success": true,
  "response_size": 1024,
  "timestamp": "2024-01-01T12:00:00Z"
}
```

#### POST `/api/load`
Inicia um teste de carga.

**Request Body:**
```json
{
  "url": "https://api.exemplo.com/endpoint",
  "requests": 100,
  "concurrency": 10
}
```

#### GET `/api/load-results`
Retorna o histórico de testes de carga.

**Query Parameters:**
- `limit`: Número máximo de resultados (padrão: 50)

#### GET `/api/load-results/{id}`
Retorna detalhes de um teste de carga específico.

## 🎯 Casos de Uso

- **Desenvolvimento de APIs**: Teste suas APIs durante o desenvolvimento
- **CI/CD**: Integre testes de performance em pipelines
- **Monitoramento**: Monitore a performance de APIs em produção
- **Análise de Performance**: Identifique gargalos e otimize sua API
- **Documentação**: Demonstre a performance da sua API

## 🏗️ Estrutura do Projeto

```
api-monitor/
├── main.go                 # Ponto de entrada da aplicação
├── go.mod                  # Dependências do Go
├── internal/
│   ├── api/               # Handlers da API REST
│   │   └── api.go
│   ├── monitor/           # Lógica de monitoramento
│   │   └── monitor.go
│   └── storage/          # Sistema de armazenamento
│       └── storage.go
└── web/
    ├── index.html        # Interface web
    └── static/
        ├── style.css     # Estilos
        └── app.js        # JavaScript
```

## 🔧 Funcionalidades Técnicas Demonstradas

- ✅ Concorrência com Goroutines
- ✅ RESTful API
- ✅ WebSocket para tempo real
- ✅ CLI Tool
- ✅ Interface Web Moderna
- ✅ Estrutura de dados eficiente
- ✅ Tratamento de erros
- ✅ Testes de carga e performance

## 📊 Métricas Coletadas

- Tempo de resposta (mínimo, máximo, médio)
- Taxa de sucesso/erro
- Códigos de status HTTP
- Tamanho das respostas
- Duração total do teste
- Estatísticas de concorrência

## 🎨 Interface

A interface foi desenvolvida com foco em:
- Design moderno e responsivo
- Experiência de usuário intuitiva
- Visualizações claras de dados
- Atualizações em tempo real
- Feedback visual imediato

## ☁️ Deploy no Railway

Este projeto está pronto para deploy no Railway! Veja o guia completo em [DEPLOY.md](DEPLOY.md).

### Resumo Rápido:

1. **Preparar repositório no GitHub:**
   ```bash
   git init
   git add .
   git commit -m "Initial commit"
   git remote add origin https://github.com/SEU_USUARIO/api-monitor.git
   git push -u origin main
   ```

2. **Deploy no Railway:**
   - Acesse [railway.app](https://railway.app)
   - Crie novo projeto
   - Conecte com GitHub
   - Selecione o repositório
   - Railway detecta automaticamente e faz o deploy!

O projeto já está configurado com:
- ✅ Dockerfile para build
- ✅ Suporte à variável PORT do Railway
- ✅ Configuração otimizada para produção

## 🤝 Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou pull requests.

## 📝 Licença

Este projeto é open source e está disponível para uso pessoal e comercial.

## 👨‍💻 Autor

Desenvolvido como projeto de portfólio demonstrando habilidades em Go, desenvolvimento web e arquitetura de software.

---

**Nota**: Este projeto foi desenvolvido para demonstrar habilidades técnicas e servir como adição ao portfólio. É funcional e pode ser usado para testes reais de APIs.

