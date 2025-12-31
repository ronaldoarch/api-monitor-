# 🚀 Guia de Deploy no Railway

Este guia vai te ajudar a fazer o deploy do API Monitor no Railway.

## 📋 Pré-requisitos

1. Conta no [Railway](https://railway.app)
2. Conta no [GitHub](https://github.com)
3. Git instalado localmente

## 🔧 Passo a Passo

### 1. Preparar o Repositório no GitHub

#### 1.1. Inicializar Git (se ainda não fez)

```bash
cd "/Users/ronaldodiasdesousa/Desktop/algo em golang"
git init
git add .
git commit -m "Initial commit: API Monitor"
```

#### 1.2. Criar Repositório no GitHub

1. Acesse [GitHub](https://github.com/new)
2. Crie um novo repositório (ex: `api-monitor`)
3. **NÃO** inicialize com README, .gitignore ou license (já temos esses arquivos)

#### 1.3. Conectar e Fazer Push

```bash
# Adicionar remote (substitua SEU_USUARIO pelo seu usuário do GitHub)
git remote add origin https://github.com/SEU_USUARIO/api-monitor.git

# Renomear branch para main (se necessário)
git branch -M main

# Fazer push
git push -u origin main
```

### 2. Deploy no Railway

#### 2.1. Criar Novo Projeto

1. Acesse [Railway Dashboard](https://railway.app/dashboard)
2. Clique em **"New Project"**
3. Selecione **"Deploy from GitHub repo"**
4. Autorize o Railway a acessar seus repositórios (se necessário)
5. Selecione o repositório `api-monitor`

#### 2.2. Configuração Automática

O Railway vai:
- Detectar automaticamente que é um projeto Go
- Usar o Dockerfile para build
- Configurar a porta automaticamente via variável `PORT`

#### 2.3. Variáveis de Ambiente (Opcional)

O projeto funciona sem variáveis de ambiente, mas você pode adicionar se necessário:

1. No Railway, vá em **Settings** → **Variables**
2. Adicione variáveis se necessário (não é obrigatório para este projeto)

#### 2.4. Deploy

O Railway vai fazer o deploy automaticamente. Você pode acompanhar os logs em tempo real.

### 3. Acessar a Aplicação

Após o deploy:

1. Railway vai gerar uma URL automática (ex: `api-monitor-production.up.railway.app`)
2. Clique em **"Generate Domain"** para criar um domínio customizado (opcional)
3. Acesse a URL no navegador

## 🔍 Verificando o Deploy

### Testar a API

```bash
# Substitua pela sua URL do Railway
curl https://seu-projeto.up.railway.app/api/test \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"url": "https://jsonplaceholder.typicode.com/posts/1"}'
```

### Acessar o Dashboard

Abra no navegador:
```
https://seu-projeto.up.railway.app
```

## 🛠️ Troubleshooting

### Problema: Build falha

**Solução:**
- Verifique os logs no Railway
- Certifique-se de que o Dockerfile está correto
- Verifique se todas as dependências estão no `go.mod`

### Problema: Aplicação não inicia

**Solução:**
- Verifique se a variável `PORT` está sendo lida corretamente
- Veja os logs no Railway para erros específicos

### Problema: Arquivos estáticos não carregam

**Solução:**
- Certifique-se de que a pasta `web/` está sendo copiada no Dockerfile
- Verifique os caminhos no código

## 📝 Estrutura de Arquivos para Deploy

Certifique-se de que estes arquivos estão no repositório:

```
api-monitor/
├── Dockerfile          ✅ Necessário
├── railway.json        ✅ Opcional (ajuda Railway)
├── .dockerignore       ✅ Recomendado
├── go.mod             ✅ Necessário
├── go.sum             ✅ Necessário
├── main.go            ✅ Necessário
├── internal/          ✅ Necessário
├── web/               ✅ Necessário
└── README.md          ✅ Recomendado
```

## 🔄 Atualizações Futuras

Para atualizar o projeto no Railway:

```bash
git add .
git commit -m "Descrição das mudanças"
git push origin main
```

O Railway vai detectar o push e fazer o redeploy automaticamente!

## 💡 Dicas

1. **Domínio Customizado**: Você pode configurar um domínio próprio no Railway
2. **Logs**: Acompanhe os logs em tempo real no dashboard do Railway
3. **Variáveis de Ambiente**: Use para configurações sensíveis
4. **Monitoramento**: Railway oferece métricas básicas de uso

## 🎉 Pronto!

Seu API Monitor está no ar! Compartilhe a URL com outros desenvolvedores ou adicione ao seu portfólio.

