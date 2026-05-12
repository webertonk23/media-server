#!/bin/bash

# Cores para o terminal
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # Sem cor

echo -e "${BLUE}Iniciando Media Server Nativamente...${NC}"

# 1. Matar processos antigos se existirem
echo "Limpando processos antigos..."
pkill -f "./media-server"
pkill -f "./worker"

# 2. Build dos binários
echo "Compilando binários..."
go build -o media-server ./cmd/server/main.go
go build -o worker ./cmd/worker/main.go

if [ $? -ne 0 ]; then
    echo "Erro na compilação. Abortando."
    exit 1
fi

# 3. Iniciar a API em background
echo -e "${GREEN}Iniciando API na porta 9000...${NC}"
./media-server > logs_api.log 2>&1 &
API_PID=$!

# 4. Iniciar os Workers em background
echo -e "${GREEN}Iniciando Scanner Worker...${NC}"
./worker scanner > logs_scanner.log 2>&1 &
SCANNER_PID=$!

echo -e "${GREEN}Iniciando Transcoder Worker...${NC}"
./worker transcoder > logs_transcoder.log 2>&1 &
TRANSCODER_PID=$!

echo -e "${BLUE}Todos os serviços estão rodando!${NC}"
echo "PIDs: API=$API_PID, Scanner=$SCANNER_PID, Transcoder=$TRANSCODER_PID"
echo "Acesse: http://localhost:9000"
echo "Use 'tail -f logs_transcoder.log' para acompanhar o processamento."

# Aguarda os processos (opcional, mas bom para manter o script vivo se rodar no foreground)
wait $API_PID $SCANNER_PID $TRANSCODER_PID
