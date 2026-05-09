# Media Server

Media server pessoal leve inspirado em Plex/Jellyfin, otimizado para Raspberry Pi 3 B e hardware N100/N150.

## 🎬 Funcionalidades Recentes

- 🎨 **Interface Cinematográfica**: Interface moderna e responsiva com tema escuro premium, micro-animações e foco em experiência visual.
- 📊 **Monitor de Transcodificação**: Acompanhe o status das conversões de vídeo em tempo real, com detecção de erros e suporte a hardware (QSV/Pi 3).
- 📜 **Logs do Sistema**: Visualize os logs unificados da API, Scanner e Transcoder diretamente pela interface web.
- 🏷️ **Qualidade de Mídia**: Detecção e exibição automática de qualidade (4K, 1080p, 720p) nos cards e detalhes.
- 📺 **Experiência de Séries**: Lista de episódios em formato de carrossel horizontal responsivo e agrupamento inteligente de temporadas.
- ⚙️ **Configuração Dinâmica**: Gerencie diretórios e intervalos de varredura sem precisar reiniciar o servidor.

## 🎯 Características Gerais

- 🎬 **Suporte a Filmes e Séries** com detecção automática e agrupamento.
- 🔍 **Metadata automática** via TMDB (Poster, Backdrop, Sinopse).
- 📺 **Streaming Inteligente** com suporte a múltiplos formatos e seek.
- 📊 **Continue Watching**: Retome de onde parou em qualquer dispositivo.
- 🔐 **Segurança e Estabilidade**: Identificadores ULID e backend em Go robusto.
- ⚡ **Extremamente Leve**: Consumo mínimo de recursos, ideal para hardware ARM.

## 🏗️ Stack

### Backend
- **Go** 1.25+
- **Fiber** (web framework)
- **GORM** (ORM) + **SQLite**
- **FFmpeg** (transcoding & metadata)

### Frontend
- **Vue 3** (Composition API)
- **TypeScript**
- **Pinia** (state management)
- **Vite** (fast build)
- **Vanilla CSS** (Custom Cinematic Design System)

### Infraestrutura
- **Docker** + Docker Compose
- **Aceleração de Hardware**: Configurado para Intel QSV e Raspberry Pi (v4l2m2m).

## 🚀 Quick Start

### Pré-requisitos
- Docker & Docker Compose

### Instalação

1. Clone o repositório e acesse a pasta:
```bash
git clone <repo-url>
cd media-server
```

2. Configure o arquivo `.env`:
```bash
PORT=9000
TMDB_API_KEY=sua_chave_tmdb
TRANSCODE_CODEC=h264_qsv # ou libx264 para software
```

3. Inicie os serviços:
```bash
docker-compose up -d
```

4. Acesse a interface em: `http://localhost:9000`

## 📁 Estrutura do Projeto

- `cmd/server`: Ponto de entrada da API.
- `cmd/worker`: Processos de scanner e transcodificação.
- `internal/api`: Handlers e rotas HTTP.
- `internal/services`: Lógica de transcode, biblioteca e metadados.
- `frontend/src`: Código fonte da interface Vue 3.
- `data/`: Banco de dados e logs persistentes.

## 🎬 Organização de Mídia

O parser identifica automaticamente a qualidade e o tipo de mídia:

- **Filmes**: `Doutor Estranho 2022 1080p.mkv` → Título: Doutor Estranho, Ano: 2022, Qualidade: 1080p.
- **Séries**: `The Boys S01E01.mp4` → Detectado como Série, Temporada 1, Episódio 1.

## 🗺️ Roadmap Atualizado

### Fase 1 - MVP ✅
- [x] Scanner e Metadata TMDB
- [x] Streaming básico e Progresso
- [x] Interface Web inicial

### Fase 2 - Séries & UI Premium ✅
- [x] Detecção de episódios (S01E01)
- [x] Interface Netflix-like (Cinematic)
- [x] Carrossel de episódios e agrupamento por série
- [x] Tag de qualidade de mídia

### Fase 3 - Monitoramento & Gestão ✅
- [x] Monitor de Transcodificação (Real-time)
- [x] Visualizador de Logs integrado
- [x] Configurações via UI (Caminhos e Scan)
- [x] Tratamento de erros de transcode

### Fase 4 - Features Avançadas 📋
- [ ] Legendas (OpenSubtitles)
- [ ] Seleção manual de qualidade (Transcode On-the-fly)
- [ ] App Mobile (Android)
- [ ] Suporte a múltiplos usuários/perfis

## 🤝 Contribuindo

Sugestões e Pull Requests são bem-vindos para melhorar este servidor!

## 📄 Licença

MIT
