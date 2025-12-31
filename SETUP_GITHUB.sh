#!/bin/bash

# Script para configurar e fazer push para GitHub
# Execute após criar o repositório no GitHub

echo "🚀 Configurando repositório Git..."

# Verificar se já é um repositório Git
if [ ! -d ".git" ]; then
    echo "Inicializando repositório Git..."
    git init
fi

# Adicionar todos os arquivos
echo "Adicionando arquivos..."
git add .

# Fazer commit inicial
echo "Fazendo commit..."
git commit -m "Initial commit: API Monitor - Ferramenta de monitoramento de APIs em Go"

# Configurar branch main
git branch -M main

echo ""
echo "✅ Repositório local configurado!"
echo ""
echo "📝 Próximos passos:"
echo "1. Copie a URL do seu repositório no GitHub"
echo "2. Execute: git remote add origin https://github.com/ronaldoarch/api-monitor.git"
echo "3. Execute: git push -u origin main"
echo ""
echo "Ou execute manualmente os comandos acima."

