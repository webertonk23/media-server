#!/bin/sh

mkdir -p data
mkdir -p cache
mkdir -p media

if [ ! -f go.mod ]; then
    echo "Inicializando projeto Go..."
    go mod init media-server
fi

if [ ! -f .air.toml ]; then
cat <<EOF > .air.toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ./cmd/server"
bin = "./tmp/main"
full_bin = "./tmp/main"
include_ext = ["go"]
exclude_dir = ["tmp", "vendor"]
EOF
fi

echo "Organizando dependências..."

go mod tidy

echo "Iniciando Air..."

air