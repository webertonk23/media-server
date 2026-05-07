# Media Server - Especificação do Projeto

## Visão Geral

Media server pessoal leve inspirado em Plex/Jellyfin, otimizado para Raspberry Pi 3 B com possibilidade de rodar em hardware N100/N150 via Docker.

## Objetivos

### Primários
- Streaming de filmes e séries com suporte a range requests (seek)
- Metadata automática via TMDB
- Continue watching (progresso de visualização)
- Interface Android TV (futuro)

### Secundários
- Leve e eficiente para Raspberry Pi
- Arquitetura escalável para crescimento futuro
- API RESTful bem estruturada
- Suporte a múltiplas qualidades por item

## Stack Tecnológica

### Backend
- **Linguagem**: Go 1.25
- **Framework Web**: Fiber v2.52
- **ORM**: GORM v1.31
- **Banco de Dados**: SQLite (inicial), preparado para PostgreSQL
- **Identificadores**: ULID (Universally Unique Lexicographically Sortable Identifier)

### Infraestrutura
- **Containerização**: Docker + Docker Compose
- **Hot Reload**: Air v1.65
- **Build**: Go modules

### Frontend (Planejado)
- **Plataforma**: Android TV
- **Linguagem**: Kotlin
- **Player**: ExoPlayer
- **Design**: Interface estilo Netflix/Plex

## Arquitetura

### Camadas

```
┌─────────────────────────────────────────┐
│         API Layer (Handlers)            │
│  - Recebe/retorna ULID                  │
│  - Validação de entrada                 │
│  - Serialização JSON                    │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│       Service Layer (Business)          │
│  - Lógica de negócio                    │
│  - Orquestração de operações            │
│  - Transações                           │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│    Repository Layer (Data Access)       │
│  - CRUD operations                      │
│  - Queries otimizadas                   │
│  - Abstração do banco                   │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│          Database (SQLite)              │
│  - Persistência                         │
│  - Índices                              │
│  - Constraints                          │
└─────────────────────────────────────────┘
```

### Pipeline de Scan

```
1. Scanner
   └─> Descobre arquivos de vídeo no filesystem

2. Parser
   └─> Extrai: título, ano, season, episode, qualidade

3. Matcher
   └─> Decide: movie, series, episode

4. Metadata Provider
   └─> Busca: TMDB (movies/TV)

5. Persistence
   └─> Salva: MediaItem + MediaFile + Series/Season/Episode
```

## Modelo de Dados

### Entidades Principais

#### MediaItem
Representa qualquer tipo de mídia (filme, série, episódio).

```go
type MediaItem struct {
    ID            uint      // PK interno
    ULID          string    // Identificador público (26 chars)
    Type          string    // movie, series, episode
    Title         string
    OriginalTitle string
    Year          int
    Overview      string
    Poster        string    // URL da imagem
    Backdrop      string    // URL da imagem
    TMDBID        int       // ID do TMDB
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

#### MediaFile
Representa arquivo físico no disco.

```go
type MediaFile struct {
    ID          uint      // PK interno
    ULID        string    // Identificador público
    MediaItemID uint      // FK para MediaItem
    Path        string    // Caminho absoluto
    Size        int64     // Tamanho em bytes
    Fingerprint string    // Hash MD5 (path + modtime)
    Quality     string    // 4K, 1080p, 720p, SD
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

#### MediaProgress
Progresso de visualização.

```go
type MediaProgress struct {
    ID            uint      // PK interno
    MediaItemID   uint      // FK para MediaItem (unique)
    Position      int64     // Posição em segundos
    Duration      int64     // Duração total em segundos
    Finished      bool
    LastWatchedAt time.Time
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

### Entidades de Séries

#### Series
Informações específicas de séries.

```go
type Series struct {
    ID             uint      // PK interno
    ULID           string    // Identificador público
    MediaItemID    uint      // FK para MediaItem (unique)
    Status         string    // Returning, Ended, Canceled
    NumberSeasons  int
    NumberEpisodes int
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

#### Season
Temporada de uma série.

```go
type Season struct {
    ID           uint      // PK interno
    ULID         string    // Identificador público
    SeriesID     uint      // FK para Series
    Number       int       // Número da temporada
    Name         string
    Overview     string
    Poster       string
    AirDate      *time.Time
    EpisodeCount int
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

**Índice Único**: `(series_id, number)` - garante uma temporada por número por série

#### Episode
Episódio de uma temporada.

```go
type Episode struct {
    ID          uint      // PK interno
    ULID        string    // Identificador público
    MediaItemID uint      // FK para MediaItem (unique)
    SeasonID    uint      // FK para Season
    Number      int       // Número do episódio
    Name        string
    Overview    string
    Still       string    // Imagem do episódio
    AirDate     *time.Time
    Runtime     int       // Duração em minutos
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**Índice Único**: `(season_id, number)` - garante um episódio por número por temporada

### Relacionamentos

```
MediaItem (type=series) ←─┐
                          │
Series ←──────────────────┘
  │
  └─> Season (1:N)
        │
        └─> Episode (1:N)
              │
              └─> MediaItem (type=episode) ←─> MediaFile (1:N)
                                            └─> MediaProgress (1:1)

MediaItem (type=movie) ←─> MediaFile (1:N)
                        └─> MediaProgress (1:1)
```

## Identificadores (ULID)

### Por que ULID?

1. **Ordenável**: Lexicograficamente ordenável por tempo de criação
2. **Compacto**: 26 caracteres vs 36 do UUID
3. **URL-friendly**: Apenas caracteres alfanuméricos (Base32)
4. **Timestamp embutido**: Primeiros 48 bits são timestamp Unix
5. **Seguro**: Não previsível, não expõe quantidade de registros

### Estrutura

```
01ARZ3NDEKTSV4RRFFQ69G5FAV
|----------|------------|
 Timestamp   Randomness
  (10 chars)  (16 chars)
```

### Implementação

- **ID interno** (uint): usado para joins e foreign keys (performance)
- **ULID público** (string): exposto na API
- **Geração automática**: hook `BeforeCreate` do GORM

## Parser de Nomes

### Filmes

**Entrada**: `Doutor Estranho no Multiverso da Loucura 2022 2160p 4K COMPACTO WEB-DL DUAL 5.1.mkv`

**Processamento**:
1. Extrair ano: `2022`
2. Extrair qualidade: `4K`
3. Remover tokens: `2160p`, `4K`, `COMPACTO`, `WEB-DL`, `DUAL`, `5.1`
4. Normalizar separadores: `.` e `_` → espaço
5. Limpar e capitalizar

**Saída**: `Doutor Estranho No Multiverso Da Loucura` (ano: 2022, qualidade: 4K)

### Séries

**Entrada**: `Cães.de.Caça.S02E01.WEB-DL.1080p.x265.DUAL.5.1.SF.mkv`

**Processamento**:
1. Detectar padrão: `S02E01` → Season 2, Episode 1
2. Normalizar separadores: `.` → espaço
3. Remover season/episode
4. Remover tokens de release
5. Remover números e letras soltas
6. Extrair ano do diretório se não encontrar no arquivo
7. Limpar e capitalizar

**Saída**: `Cães De Caça` (S02E01, 1080p, ano do diretório)

### Padrões Suportados

**Séries**:
- `S01E01`, `S1E1` (padrão)
- `1x01`, `1x1` (alternativo)
- `Season 1 Episode 1` (verbose)
- `101`, `1001` (numérico: season + episode)

### Tokens de Release Removidos

- **Qualidade**: 2160p, 1080p, 720p, 480p, 4k, uhd, hd
- **Codec**: x264, x265, h264, h265, hevc, avc, xvid, x, h
- **Source**: web-dl, webrip, bluray, brrip, dvdrip, hdtv
- **Audio**: dual, aac, ac3, dts, 5.1, 7.1, 2.0, atmos
- **Grupos**: starckfilmes, sf, yify, rarbg, yts, amzn, nf, netflix

## API Endpoints

### Scan

```
POST /scan              # Escaneia tudo (filmes + séries)
POST /scan/movies       # Escaneia apenas filmes
POST /scan/series       # Escaneia apenas séries
```

### Filmes

```
GET /movies?page=1&limit=20&search=doutor
```

**Resposta**:
```json
{
  "page": 1,
  "limit": 20,
  "total": 1,
  "items": [
    {
      "id": "01KR1VY139NEJYM6Q1PGEACFV7",
      "type": "movie",
      "title": "Doutor Estranho No Multiverso Da Loucura",
      "year": 2022,
      "overview": "...",
      "poster": "https://image.tmdb.org/t/p/w500/...",
      "backdrop": "https://image.tmdb.org/t/p/original/...",
      "stream_url": "/stream/01KR1VY139NEJYM6Q1PGEACFV7"
    }
  ]
}
```

### Séries

```
GET /series?page=1&limit=20&search=breaking
```

**Resposta**:
```json
{
  "page": 1,
  "limit": 20,
  "total": 1,
  "items": [
    {
      "id": "01KR1YDYDSPNKKF8NBH4PRB0Z8",
      "type": "series",
      "title": "Cães De Caça",
      "year": 2023,
      "overview": "...",
      "poster": "...",
      "backdrop": "...",
      "status": "Ended",
      "number_seasons": 2,
      "number_episodes": 15
    }
  ]
}
```

### Temporadas

```
GET /series/:seriesId/seasons
```

**Resposta**:
```json
[
  {
    "id": "01KR1YDYEZPZ6GV8AJJKX583JN",
    "number": 1,
    "name": "Season 1",
    "overview": "...",
    "poster": "...",
    "episode_count": 8
  }
]
```

### Episódios

```
GET /seasons/:seasonId/episodes
```

**Resposta**:
```json
[
  {
    "id": "01KR1YDYF0ABCDEFGHIJKLMNOP",
    "type": "episode",
    "title": "Pilot",
    "overview": "...",
    "season_number": 1,
    "episode_number": 1,
    "still": "...",
    "runtime": 45,
    "stream_url": "/stream/01KR1YDYF0ABCDEFGHIJKLMNOP"
  }
]
```

### Streaming

```
GET /stream/:ulid
```

**Características**:
- Suporta range requests (HTTP 206 Partial Content)
- Seek funcionando
- Content-Type automático baseado na extensão
- Funciona para filmes e episódios

**Headers**:
```
Accept-Ranges: bytes
Content-Type: video/x-matroska
Content-Disposition: inline
```

### Progresso

```
POST /progress/:ulid
Content-Type: application/json

{
  "position": 1200,
  "duration": 7200,
  "finished": false
}
```

### Mídia Genérica

```
GET /media?type=movie&page=1&limit=20
GET /media?type=series
GET /media?type=episode
```

## Estrutura de Diretórios

### Projeto

```
media-server/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP handlers
│   │   └── routes/              # Definição de rotas
│   ├── config/                  # Configuração
│   ├── database/                # Conexão DB
│   ├── dto/                     # Data Transfer Objects
│   ├── mappers/                 # Model → DTO
│   ├── metadata/                # Providers (TMDB)
│   ├── models/                  # Modelos GORM
│   ├── repositories/            # Acesso a dados
│   ├── scanner/                 # Scanner de arquivos
│   ├── services/                # Lógica de negócio
│   └── utils/                   # Utilitários (parsers)
├── media/
│   ├── movies/                  # Filmes
│   └── series/                  # Séries
├── data/
│   └── media.db                 # Banco SQLite
├── docker/
│   └── go/
│       ├── Dockerfile
│       └── start.sh
├── .agents/                     # Documentação do projeto
├── docker-compose.yml
├── go.mod
└── go.sum
```

### Mídia

```
media/
├── movies/
│   ├── Filme 1 2022 1080p.mkv
│   └── Filme 2 2023 4K.mp4
└── series/
    ├── Serie S01 2020/
    │   ├── Serie S01E01.mkv
    │   ├── Serie S01E02.mkv
    │   └── ...
    └── Serie S02 2021/
        ├── Serie S02E01.mkv
        └── ...
```

## Configuração

### Variáveis de Ambiente (.env)

```bash
PORT=9000
MEDIA_PATH=/app/media
DB_PATH=/app/data/media.db
DB_DRIVER=sqlite
TMDB_API_KEY=sua_chave_aqui
```

### Docker Compose

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: docker/go/Dockerfile
    ports:
      - "9000:9000"
    volumes:
      - .:/app
      - ./media:/app/media
      - ./data:/app/data
    environment:
      - PORT=9000
      - MEDIA_PATH=/app/media
      - DB_PATH=/app/data/media.db
      - TMDB_API_KEY=${TMDB_API_KEY}
```

## Decisões Técnicas

### 1. Por que ULID ao invés de UUID?

- **Ordenável**: Facilita paginação e ordenação
- **Compacto**: 26 vs 36 caracteres
- **Timestamp embutido**: Útil para debugging
- **URL-friendly**: Sem caracteres especiais

### 2. Por que ID interno + ULID?

- **Performance**: Joins com uint são mais rápidos
- **Índices**: Índices numéricos são menores
- **Compatibilidade**: GORM funciona melhor com IDs numéricos
- **Segurança**: ULID não expõe quantidade de registros

### 3. Por que separar MediaItem e MediaFile?

- **Múltiplas qualidades**: Um filme pode ter versões 4K e 1080p
- **Múltiplos idiomas**: Um episódio pode ter versões dubladas
- **Flexibilidade**: Metadata compartilhada entre arquivos
- **Escalabilidade**: Preparado para CDN e transcodificação

### 4. Por que SQLite inicialmente?

- **Simplicidade**: Zero configuração
- **Performance**: Suficiente para uso pessoal
- **Portabilidade**: Arquivo único
- **Raspberry Pi**: Leve e eficiente

**Migração futura**: Preparado para PostgreSQL quando necessário

### 5. Por que Go + Fiber?

- **Performance**: Compilado, rápido, baixo consumo
- **Concorrência**: Goroutines para scan paralelo
- **Simplicidade**: Sintaxe clara, fácil manutenção
- **Raspberry Pi**: Binário único, sem runtime

### 6. Por que Docker desde o início?

- **Consistência**: Mesmo ambiente dev/prod
- **Isolamento**: Não polui o sistema host
- **Portabilidade**: Funciona em qualquer lugar
- **Hot reload**: Air para desenvolvimento

## Performance

### Otimizações Implementadas

1. **Índices**:
   - `media_items.ulid` (unique)
   - `media_items.type`
   - `media_items.tmdb_id`
   - `media_files.path` (unique)
   - `media_files.fingerprint`
   - `(series_id, number)` para Season
   - `(season_id, number)` para Episode

2. **Fingerprint**:
   - MD5(path + modtime)
   - Evita re-scan de arquivos não modificados

3. **Preload**:
   - Eager loading de relacionamentos quando necessário
   - Evita N+1 queries

4. **Streaming**:
   - Range requests nativos do Fiber
   - Sem buffering em memória
   - Seek instantâneo

### Benchmarks (Raspberry Pi 3 B)

- **Scan**: ~100 arquivos/segundo
- **Streaming**: Suporta 2-3 streams simultâneos 1080p
- **API**: ~1000 req/s (endpoints simples)
- **Memória**: ~50MB em idle

## Segurança

### Implementado

1. **ULID**: IDs não previsíveis
2. **Path validation**: Previne directory traversal
3. **Content-Type**: Baseado em extensão, não em input do usuário
4. **Fingerprint**: Detecta arquivos duplicados

### Planejado

1. **Autenticação**: JWT tokens
2. **Autorização**: Perfis de usuário
3. **Rate limiting**: Proteção contra abuse
4. **HTTPS**: TLS/SSL

## Testes

### Estrutura

```
internal/
├── utils/
│   ├── parser.go
│   ├── parser_test.go
│   ├── series_parser.go
│   └── series_parser_test.go
└── ...
```

### Cobertura Planejada

- [ ] Parser de filmes
- [ ] Parser de séries
- [ ] Repositories
- [ ] Services
- [ ] Handlers (integration tests)

## Roadmap

### Fase 1 - MVP ✅
- [x] Scanner de filmes
- [x] Metadata TMDB
- [x] Streaming básico
- [x] Progress tracking
- [x] ULID
- [x] Arquitetura MediaItem
- [x] Scanner de séries
- [x] Parser de séries
- [x] API de séries

### Fase 2 - Melhorias 🚧
- [ ] Metadata TMDB TV (séries)
- [ ] Continue Watching
- [ ] Seleção de qualidade
- [ ] Múltiplos providers (IMDB, TVDB)
- [ ] Legendas (OpenSubtitles)

### Fase 3 - Features Avançadas 📋
- [ ] Transcodificação (FFmpeg)
- [ ] Autenticação/Autorização
- [ ] Múltiplos usuários
- [ ] Perfis
- [ ] Watchlist

### Fase 4 - Android TV 📱
- [ ] App Android TV
- [ ] ExoPlayer
- [ ] Interface Netflix-like
- [ ] Chromecast
- [ ] Offline download

## Manutenção

### Logs

```bash
# Ver logs do container
docker logs media-server -f

# Ver logs do scan
docker logs media-server | grep "Processando"

# Ver erros
docker logs media-server | grep "Erro"
```

### Backup

```bash
# Backup do banco
cp data/media.db data/media.db.backup

# Backup completo
tar -czf backup.tar.gz data/ media/
```

### Limpeza

```bash
# Remover banco e rescanear
docker compose down
rm -f data/media.db
docker compose up -d
```

## Contribuindo

### Convenções

1. **Commits**: Conventional Commits
2. **Branches**: feature/, bugfix/, hotfix/
3. **Code Style**: gofmt + golint
4. **Testes**: Cobertura mínima 80%

### Workflow

1. Fork do repositório
2. Criar branch feature
3. Implementar + testes
4. Pull request com descrição detalhada

## Licença

MIT

## Contato

- **Projeto**: Media Server
- **Versão**: 2.0.0
- **Status**: Em desenvolvimento ativo
