# 📚 Exemplos de Uso

## Exemplos Práticos do API Monitor

### 1. Teste Rápido via CLI

```bash
# Teste uma API pública
./api-monitor --cli test https://jsonplaceholder.typicode.com/posts/1

# Ou usando go run
go run main.go --cli test https://httpbin.org/get
```

### 2. Teste de Carga via CLI

```bash
# 100 requisições com concorrência de 10
./api-monitor --cli load https://jsonplaceholder.typicode.com/posts 100 10

# 500 requisições com concorrência de 50 (teste mais intenso)
./api-monitor --cli load https://httpbin.org/delay/1 500 50
```

### 3. Usando a API REST

#### Teste Rápido via cURL

```bash
curl -X POST http://localhost:8080/api/test \
  -H "Content-Type: application/json" \
  -d '{"url": "https://jsonplaceholder.typicode.com/posts/1"}'
```

#### Teste de Carga via cURL

```bash
curl -X POST http://localhost:8080/api/load \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://httpbin.org/get",
    "requests": 100,
    "concurrency": 10
  }'
```

#### Obter Histórico

```bash
# Listar últimos 20 testes
curl http://localhost:8080/api/load-results?limit=20

# Obter detalhes de um teste específico
curl http://localhost:8080/api/load-results/{test-id}
```

### 4. Exemplo em JavaScript (Frontend)

```javascript
// Teste rápido
async function runQuickTest(url) {
  const response = await fetch('/api/test', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ url }),
  });
  
  const result = await response.json();
  console.log('Resultado:', result);
  return result;
}

// Teste de carga
async function runLoadTest(url, requests, concurrency) {
  const response = await fetch('/api/load', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ url, requests, concurrency }),
  });
  
  const result = await response.json();
  console.log('Teste iniciado:', result);
  return result;
}
```

### 5. Exemplo em Python

```python
import requests
import json

# Teste rápido
def test_api(url):
    response = requests.post(
        'http://localhost:8080/api/test',
        json={'url': url}
    )
    return response.json()

# Teste de carga
def load_test(url, requests=100, concurrency=10):
    response = requests.post(
        'http://localhost:8080/api/load',
        json={
            'url': url,
            'requests': requests,
            'concurrency': concurrency
        }
    )
    return response.json()

# Exemplo de uso
result = test_api('https://jsonplaceholder.typicode.com/posts/1')
print(json.dumps(result, indent=2))
```

### 6. Integração com CI/CD (GitHub Actions)

```yaml
name: API Performance Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21
      
      - name: Build API Monitor
        run: go build -o api-monitor
      
      - name: Run Performance Test
        run: |
          ./api-monitor --cli load https://api.exemplo.com/endpoint 100 10
```

### 7. APIs Públicas para Testes

Aqui estão algumas APIs públicas que você pode usar para testar:

- **JSONPlaceholder**: `https://jsonplaceholder.typicode.com/posts/1`
- **HTTPBin**: `https://httpbin.org/get`
- **ReqRes**: `https://reqres.in/api/users/1`
- **Random User**: `https://randomuser.me/api/`

### 8. Monitoramento em Tempo Real

O dashboard web suporta atualizações em tempo real via WebSocket. Quando você executa um teste através da interface web, os resultados aparecem automaticamente sem necessidade de atualizar a página.

### 9. Script de Automação (Bash)

```bash
#!/bin/bash

# Script para executar testes periódicos
API_URL="https://api.exemplo.com/health"
MONITOR="./api-monitor"

echo "Executando teste de saúde da API..."
$MONITOR --cli test $API_URL

echo "Executando teste de carga..."
$MONITOR --cli load $API_URL 50 5

echo "Testes concluídos!"
```

### 10. Análise de Resultados

Os resultados incluem:
- Tempo de resposta (mínimo, máximo, médio)
- Taxa de sucesso/erro
- Distribuição de códigos de status HTTP
- Tamanho das respostas
- Duração total do teste

Use essas métricas para:
- Identificar gargalos de performance
- Validar SLAs de resposta
- Comparar performance entre versões
- Detectar problemas de escalabilidade

