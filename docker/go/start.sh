#!/bin/sh

# Garantir pastas base
mkdir -p data cache media

echo "--- Preparando ambiente Go ---"
go mod tidy

echo "--- Buildando Frontend ---"
cd frontend
npm install  
npm run build
cd ..

echo "--- Compilando binários Go ---"
go build -o ./tmp/api ./cmd/server
go build -o ./tmp/worker ./cmd/worker

echo "--- Iniciando API na porta 9000 ---"
./tmp/api &
API_PID=$!

sleep 5

echo "--- Iniciando Serviços Secundários ---"

./tmp/worker scanner &
SCANNER_PID=$!
echo "[PID $SCANNER_PID] Scanner iniciado"

#./tmp/worker transcoder &
#TRANSCODER_PID=$!
#echo "[PID $TRANSCODER_PID] Transcoder iniciado"

echo "--- Todos os serviços iniciados ---"
wait $API_PID