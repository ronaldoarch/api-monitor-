let ws = null;
let currentLoadTestId = null;

// Conectar WebSocket
function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(`${protocol}//${window.location.host}/ws`);
    
    ws.onmessage = function(event) {
        const message = JSON.parse(event.data);
        handleWebSocketMessage(message);
    };
    
    ws.onerror = function(error) {
        console.error('WebSocket error:', error);
    };
    
    ws.onclose = function() {
        setTimeout(connectWebSocket, 3000);
    };
}

function handleWebSocketMessage(message) {
    if (message.type === 'load_test_result') {
        displayLoadTestResult(message.data);
    } else if (message.type === 'test_result') {
        displayQuickTestResult(message.data);
    }
}

// Tabs
function switchTab(tabName) {
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.classList.remove('active');
    });
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    
    document.getElementById(`${tabName}-tab`).classList.add('active');
    event.target.classList.add('active');
    
    if (tabName === 'history') {
        loadHistory();
    }
}

// Teste Rápido
document.getElementById('quick-test-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const url = document.getElementById('quick-url').value;
    const resultContainer = document.getElementById('quick-result');
    
    resultContainer.innerHTML = '<div class="loading"><div class="spinner"></div>Executando teste...</div>';
    
    try {
        const response = await fetch('/api/test', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ url }),
        });
        
        const result = await response.json();
        displayQuickTestResult(result);
    } catch (error) {
        resultContainer.innerHTML = `<div class="result-card error"><h3>Erro</h3><p>${error.message}</p></div>`;
    }
});

function displayQuickTestResult(result) {
    const resultContainer = document.getElementById('quick-result');
    const success = result.success;
    const statusClass = success ? 'success' : 'error';
    
    resultContainer.innerHTML = `
        <div class="result-card ${statusClass}">
            <h3>Resultado do Teste</h3>
            <div class="metrics">
                <div class="metric">
                    <div class="metric-label">Status</div>
                    <div class="metric-value ${success ? 'success' : 'error'}">${result.status || 'N/A'}</div>
                </div>
                <div class="metric">
                    <div class="metric-label">Tempo de Resposta</div>
                    <div class="metric-value">${result.duration}ms</div>
                </div>
                <div class="metric">
                    <div class="metric-label">Tamanho da Resposta</div>
                    <div class="metric-value">${formatBytes(result.response_size || 0)}</div>
                </div>
                <div class="metric">
                    <div class="metric-label">Sucesso</div>
                    <div class="metric-value ${success ? 'success' : 'error'}">${success ? 'Sim' : 'Não'}</div>
                </div>
            </div>
            ${result.error ? `<p style="margin-top: 15px; color: var(--danger);"><strong>Erro:</strong> ${result.error}</p>` : ''}
            <p style="margin-top: 15px; color: var(--gray); font-size: 0.875rem;">
                <strong>URL:</strong> ${result.url}<br>
                <strong>Timestamp:</strong> ${new Date(result.timestamp).toLocaleString('pt-BR')}
            </p>
        </div>
    `;
}

// Teste de Carga
document.getElementById('load-test-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const url = document.getElementById('load-url').value;
    const requests = parseInt(document.getElementById('requests').value);
    const concurrency = parseInt(document.getElementById('concurrency').value);
    const resultContainer = document.getElementById('load-result');
    
    resultContainer.innerHTML = '<div class="loading"><div class="spinner"></div>Iniciando teste de carga...</div>';
    
    try {
        const response = await fetch('/api/load', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ url, requests, concurrency }),
        });
        
        const result = await response.json();
        resultContainer.innerHTML = `
            <div class="result-card">
                <h3>Teste de Carga Iniciado</h3>
                <p>O teste está sendo executado. Os resultados aparecerão aqui quando concluído.</p>
                <div class="loading"><div class="spinner"></div>Executando ${requests} requisições com concorrência de ${concurrency}...</div>
            </div>
        `;
    } catch (error) {
        resultContainer.innerHTML = `<div class="result-card error"><h3>Erro</h3><p>${error.message}</p></div>`;
    }
});

function displayLoadTestResult(result) {
    const resultContainer = document.getElementById('load-result');
    const successRate = ((result.success_count / result.total_requests) * 100).toFixed(2);
    
    resultContainer.innerHTML = `
        <div class="result-card">
            <h3>Resultado do Teste de Carga</h3>
            <div class="metrics">
                <div class="metric">
                    <div class="metric-label">Total de Requisições</div>
                    <div class="metric-value">${result.total_requests}</div>
                </div>
                <div class="metric">
                    <div class="metric-label">Taxa de Sucesso</div>
                    <div class="metric-value success">${successRate}%</div>
                </div>
                <div class="metric">
                    <div class="metric-label">Requisições Bem-sucedidas</div>
                    <div class="metric-value success">${result.success_count}</div>
                </div>
                <div class="metric">
                    <div class="metric-label">Requisições com Erro</div>
                    <div class="metric-value error">${result.error_count}</div>
                </div>
                <div class="metric">
                    <div class="metric-label">Tempo Médio de Resposta</div>
                    <div class="metric-value">${result.avg_response_time.toFixed(2)}ms</div>
                </div>
                <div class="metric">
                    <div class="metric-label">Tempo Mínimo</div>
                    <div class="metric-value">${result.min_response_time}ms</div>
                </div>
                <div class="metric">
                    <div class="metric-label">Tempo Máximo</div>
                    <div class="metric-value">${result.max_response_time}ms</div>
                </div>
                <div class="metric">
                    <div class="metric-label">Duração Total</div>
                    <div class="metric-value">${result.duration}ms</div>
                </div>
            </div>
            <div style="margin-top: 20px;">
                <h4 style="margin-bottom: 10px;">Códigos de Status HTTP:</h4>
                <div style="display: flex; flex-wrap: wrap; gap: 10px;">
                    ${Object.entries(result.status_codes).map(([code, count]) => 
                        `<span class="status-badge ${code < 400 ? 'success' : 'error'}">${code}: ${count}</span>`
                    ).join('')}
                </div>
            </div>
            <p style="margin-top: 15px; color: var(--gray); font-size: 0.875rem;">
                <strong>URL:</strong> ${result.url}<br>
                <strong>Concorrência:</strong> ${result.concurrency}<br>
                <strong>Timestamp:</strong> ${new Date(result.timestamp).toLocaleString('pt-BR')}
            </p>
        </div>
    `;
}

// Histórico
async function loadHistory() {
    const historyContent = document.getElementById('history-content');
    historyContent.innerHTML = '<div class="loading"><div class="spinner"></div>Carregando histórico...</div>';
    
    try {
        const response = await fetch('/api/load-results?limit=20');
        const results = await response.json();
        
        if (results.length === 0) {
            historyContent.innerHTML = '<p style="text-align: center; color: var(--gray); padding: 40px;">Nenhum teste encontrado no histórico.</p>';
            return;
        }
        
        historyContent.innerHTML = results.map(result => `
            <div class="history-item" onclick="viewTestDetails('${result.id}')">
                <div class="history-item-header">
                    <div class="history-item-url">${result.url}</div>
                    <span class="status-badge ${result.success_count > result.error_count ? 'success' : 'error'}">
                        ${((result.success_count / result.total_requests) * 100).toFixed(1)}% sucesso
                    </span>
                </div>
                <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 10px; margin-top: 10px; font-size: 0.875rem;">
                    <div><strong>Requisições:</strong> ${result.total_requests}</div>
                    <div><strong>Tempo Médio:</strong> ${result.avg_response_time.toFixed(2)}ms</div>
                    <div><strong>Duração:</strong> ${result.duration}ms</div>
                </div>
                <div class="history-item-time">${new Date(result.timestamp).toLocaleString('pt-BR')}</div>
            </div>
        `).join('');
    } catch (error) {
        historyContent.innerHTML = `<div class="result-card error"><h3>Erro</h3><p>${error.message}</p></div>`;
    }
}

async function viewTestDetails(id) {
    try {
        const response = await fetch(`/api/load-results/${id}`);
        const result = await response.json();
        displayLoadTestResult(result);
        switchTab('load');
        document.querySelectorAll('.tab-btn')[1].click();
    } catch (error) {
        alert('Erro ao carregar detalhes do teste');
    }
}

function formatBytes(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

// Inicializar
connectWebSocket();

