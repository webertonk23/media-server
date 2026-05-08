#!/bin/sh

# Garantir pastas base
mkdir -p data cache media

echo "--- Preparando ambiente Go ---"
go mod tidy

echo "--- Buildando Frontend ---"
# Entra na pasta, instala dependências (opcional se já estiver no container) e builda
cd frontend
npm install  # Remova essa linha se o install já for feito no Dockerfile
npm run build
cd ..

echo "--- Compilando binários Go ---"
go build -o ./tmp/api ./cmd/server
go build -o ./tmp/worker ./cmd/worker

echo "--- Iniciando Serviços (PIDs Separados) ---"

# 1. Inicia o Scanner em background
./tmp/worker scanner &
SCANNER_PID=$!
echo "[PID $SCANNER_PID] Scanner iniciado"

# 2. Inicia o Transcoder em background
./tmp/worker transcoder &
TRANSCODER_PID=$!
echo "[PID $TRANSCODER_PID] Transcoder iniciado"

# 3. Inicia a API (Processo principal que segura o container vivo)
echo "--- Iniciando API na porta 9000 ---"
exec ./tmp/api