# Media Server

Media server pessoal leve inspirado em Plex/Jellyfin, otimizado para Raspberry Pi 3 B e hardware N100/N150.

## 🎯 Características

- 🎬 **Suporte a Filmes e Séries** com detecção automática
- 🔍 **Metadata automática** via TMDB
- 📺 **Streaming** com range requests e suporte a múltiplos formatos
- 📊 **Continue Watching** (progresso de visualização por episódio)
- 🔐 **IDs seguros** com ULID
- 🎨 **Interface Web moderna** com Vue 3 e Tailwind CSS
- 🐳 **Docker** desde o início
- ⚡ **Leve e rápido** para Raspberry Pi

## 🏗️ Stack

### Backend
- **Go** 1.23+
- **Fiber** (web framework)
- **GORM** (ORM)
- **SQLite** (banco de dados)
- **ULID** (identificadores únicos)

### Frontend
- **Vue 3** (SPA)
- **TypeScript**
- **Tailwind CSS** (styling)
- **Vite** (build tool)
- **Pinia** (state management)

### Infraestrutura
- **Docker** + Docker Compose
- **Air** (hot reload Go)

## 🚀 Quick Start

### Pré-requisitos
- Docker
- Docker Compose

### Instalação

1. Clone o repositório:
```bash
git clone <repo-url>
cd media-server
```

2. Configure o `.env`:
```bash
PORT=9000
MEDIA_PATH=/app/media
DB_PATH=/app/data/media.db
TMDB_API_KEY=sua_chave_aqui
```

3. Inicie o servidor:
```bash
docker-compose up
```

4. Acesse:
```
API:       http://localhost:9000/health
Interface: http://localhost:5173 (desenvolvimento)
```

## 📁 Estrutura de Diretórios

```
media-server/
├── cmd/
│   └── server/          # Entry point
├── internal/
│   ├── api/
│   │   ├── handlers/    # HTTP handlers
│   │   └── routes/      # Rotas
│   ├── config/          # Configuração
│   ├── database/        # Conexão DB
│   ├── dto/             # Data Transfer Objects
│   ├── mappers/         # Model → DTO
│   ├── metadata/        # Providers (TMDB)
│   ├── models/          # Modelos GORM
│   ├── repositories/    # Acesso a dados
│   ├── scanner/         # Scanner de arquivos
│   ├── services/        # Lógica de negócio
│   └── utils/           # Utilitários
├── media/
│   └── movies/          # Seus filmes aqui
├── data/
│   └── media.db         # Banco SQLite
└── docker/
    └── go/              # Dockerfile + scripts
```

### Frontend

```
frontend/
├── src/
│   ├── components/        # Componentes reutilizáveis
│   │   ├── common/        # Componentes genéricos
│   │   ├── media/         # Componentes de mídia
│   │   ├── player/        # Componentes do player
│   │   └── search/        # Componentes de busca
│   ├── composables/       # Composições Vue 3
│   │   ├── useVideoPlayer.ts
│   │   ├── useInfiniteScroll.ts
│   │   └── ...
│   ├── pages/             # Páginas da aplicação
│   │   ├── HomePage.vue
│   │   ├── CatalogPage.vue
│   │   ├── PlayerPage.vue
│   │   └── ...
│   ├── router/            # Configuração de rotas
│   ├── stores/            # Pinia stores (state)
│   │   ├── mediaStore.ts
│   │   ├── playerStore.ts
│   │   └── ...
│   ├── services/          # Serviços (API calls, etc)
│   ├── types/             # TypeScript types
│   └── utils/             # Utilitários
├── vite.config.ts
├── tailwind.config.js
└── package.json
```

## 🎬 Organização de Mídia

### Filmes
Coloque seus filmes em `media/movies/`:

```
media/movies/
├── Doutor Estranho no Multiverso da Loucura 2022 2160p 4K WEB-DL DUAL 5.1.mkv
├── Avatar O Caminho da Água 2022 1080p BluRay.mp4
└── Top Gun Maverick 2022 720p WEB-DL.mkv
```

### Séries
Coloque suas séries em `media/series/`:

```
media/series/
├── Breaking Bad S01/
│   ├── Breaking.Bad.S01E01.Pilot.1080p.mkv
│   ├── Breaking.Bad.S01E02.Cat's.in.the.Bag.1080p.mkv
│   └── ...
└── Game of Thrones S01/
    ├── Game.of.Thrones.S01E01.Winter.is.Coming.mkv
    └── ...
```

O parser extrai automaticamente:
- **Título**: "Doutor Estranho No Multiverso Da Loucura"
- **Ano**: 2022
- **Qualidade**: 4K, 1080p, 720p
- **Série/Temporada/Episódio**: S01E01, S02E05, etc.

## 📡 API Endpoints

### Health Check
```bash
GET /health
```

### Listar Filmes
```bash
GET /movies?page=1&limit=20&search=doutor
```

Resposta:
```json
{
  "page": 1,
  "limit": 20,
  "total": 1,
  "items": [
    {
      "id": "01JH8X9K2MZQW3R5T7V9Y1B4D6",
      "type": "movie",
      "title": "Doutor Estranho No Multiverso Da Loucura",
      "year": 2022,
      "overview": "...",
      "poster": "https://image.tmdb.org/t/p/w500/...",
      "backdrop": "https://image.tmdb.org/t/p/original/...",
      "stream_url": "/stream/01JH8X9K2MZQW3R5T7V9Y1B4D6"
    }
  ]
}
```

### Listar Todos os Itens
```bash
GET /media?type=movie&page=1&limit=20
```

### Streaming
```bash
GET /stream/:ulid
```

Suporta:
- Range requests (seek)
- Múltiplos formatos: MP4, MKV, AVI, MOV, WEBM

### Atualizar Progresso
```bash
POST /progress/:ulid
Content-Type: application/json

{
  "position": 1200,
  "duration": 7200,
  "finished": false
}
```

## 🔍 Scanner

O scanner percorre automaticamente o diretório de mídia e:

1. **Descobre** arquivos de vídeo
2. **Extrai** informações do nome (título, ano, qualidade)
3. **Busca** metadata no TMDB
4. **Salva** no banco de dados

### Executar Scan Manualmente

```bash
# Dentro do container
docker exec -it media-server-app-1 /bin/sh
go run cmd/server/main.go scan
```

Ou adicione um endpoint:
```go
app.Post("/scan", func(c *fiber.Ctx) error {
    service := services.NewLibraryService()
    err := service.ScanMovies()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"success": true})
})
```

## 🗄️ Banco de Dados

### Modelos

#### MediaItem
Representa qualquer tipo de mídia (filme, série, episódio):
```go
type MediaItem struct {
    ID            uint
    ULID          string  // Identificador público
    Type          string  // movie, series, episode
    Title         string
    OriginalTitle string
    Year          int
    Overview      string
    Poster        string
    Backdrop      string
    TMDBID        int
}
```

#### MediaFile
Representa arquivo físico:
```go
type MediaFile struct {
    ID          uint
    ULID        string
    MediaItemID uint    // FK para MediaItem
    Path        string
    Size        int64
    Fingerprint string
    Quality     string  // 4K, 1080p, 720p, SD
}
```

#### MediaProgress
Progresso de visualização:
```go
type MediaProgress struct {
    ID            uint
    MediaItemID   uint
    Position      int64
    Duration      int64
    Finished      bool
    LastWatchedAt time.Time
}
```

## 🎨 Arquitetura

### Pipeline de Scan

```
Scanner → Parser → Matcher → Metadata → Persistence
   ↓         ↓        ↓          ↓           ↓
Arquivos  Título   Tipo      TMDB      MediaItem
                                       + MediaFile
```

### Camadas

```
┌─────────────────────────────────────────┐
│           API Layer (Handlers)          │
│  - Recebe/retorna ULID                  │
│  - Validação de entrada                 │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│         Service Layer (Business)        │
│  - Lógica de negócio                    │
│  - Orquestração                         │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│      Repository Layer (Data Access)     │
│  - CRUD operations                      │
│  - Queries                              │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│            Database (SQLite)            │
└─────────────────────────────────────────┘
```

### Frontend - Arquitetura

```
┌─────────────────────────────────────────┐
│         Vue Router (Pages)              │
│  - HomePage, CatalogPage, PlayerPage    │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│    Components (UI renderização)         │
│  - Reutilizáveis em toda aplicação      │
│  - Media, Player, Search, Common        │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│       Pinia Stores (State Management)   │
│  - mediaStore, playerStore, uiStore     │
│  - Gerencia estado global                │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│   Services (API Integration)            │
│  - api.ts, mediaService, streamService  │
│  - Comunicação com backend               │
└─────────────────────────────────────────┘
```

### Composables

Hooks Vue 3 reutilizáveis:
- `useVideoPlayer` - Controle do player
- `useInfiniteScroll` - Carregamento infinito
- `useKeyboardShortcuts` - Atalhos de teclado
- `useMediaQuery` - Media queries responsivas
- `useDebounce` - Debounce para busca

```

### ULID
- IDs não previsíveis
- Não expõe quantidade de registros
- Ordenáveis por tempo
- URL-friendly

Exemplo: `01JH8X9K2MZQW3R5T7V9Y1B4D6`

## ⚡ Performance

### Otimizações
- ID interno (uint) para joins
- Índices em campos críticos
- Foreign keys eficientes
- ULID indexado

### Índices Criados
- `media_items.ulid` (unique)
- `media_items.type`
- `media_files.path` (unique)
- `media_files.fingerprint`
- `media_progress.media_item_id` (unique)

## 🧪 Testes

```bash
# Compilar
go build -o tmp/test ./cmd/server

# Testar health
curl http://localhost:9000/health

# Listar filmes
curl http://localhost:9000/movies

# Stream (substitua pelo ULID real)
curl http://localhost:9000/stream/01JH8X9K2MZQW3R5T7V9Y1B4D6 \
  -H "Range: bytes=0-1024"
```

## 🛠️ Desenvolvimento

### Backend - Hot Reload
O projeto usa Air para hot reload automático:

```bash
docker-compose up
# Edite qualquer arquivo .go
# O servidor reinicia automaticamente
```

### Frontend - Desenvolvimento
```bash
cd frontend
npm install
npm run dev
# Acesse http://localhost:5173
```

Build para produção:
```bash
npm run build
```

### Adicionar Dependências (Backend)
```bash
docker exec -it media-server-app-1 /bin/sh
go get github.com/alguma/lib
```

### Adicionar Dependências (Frontend)
```bash
cd frontend
npm install nome-do-pacote
```

## 📚 Documentação

- [Frontend README](frontend/README.md) - Informações sobre a interface Vue 3
- [PERFORMANCE_OPTIMIZATIONS.md](frontend/PERFORMANCE_OPTIMIZATIONS.md) - Otimizações do frontend
- [REFACTORING.md](REFACTORING.md) - Detalhes da refatoração
- [ULID_MIGRATION.md](ULID_MIGRATION.md) - Migração para ULID
- [CHANGELOG.md](CHANGELOG.md) - Histórico de mudanças

## 🗺️ Roadmap

### Fase 1 - MVP ✅
- [x] Scanner de filmes
- [x] Metadata TMDB
- [x] Streaming básico
- [x] Progress tracking
- [x] ULID
- [x] Arquitetura MediaItem
- [x] Interface Web Vue 3

### Fase 2 - Séries 🚧
- [x] Detecção de séries (S01E01)
- [x] Modelos: Series, Season, Episode
- [x] Scanner de séries
- [ ] Metadata de séries (melhorias)
- [ ] Interface de séries (refinamentos)

### Fase 3 - Features 📋
- [x] Continue Watching
- [ ] Seleção de qualidade
- [ ] Múltiplos providers (IMDB, TVDB)
- [ ] Legendas (OpenSubtitles)
- [ ] Transcodificação (FFmpeg)
- [ ] Busca avançada

### Fase 4 - Clientes Adicionais 📱
- [ ] App Android TV
- [ ] ExoPlayer
- [ ] Interface Netflix-like (mobile)
- [ ] Chromecast

## 🤝 Contribuindo

Este é um projeto pessoal, mas sugestões são bem-vindas!

## 📄 Licença

MIT

## 🙏 Créditos

- Metadata: [TMDB](https://www.themoviedb.org/)
- Inspiração: Plex, Jellyfin, Emby
