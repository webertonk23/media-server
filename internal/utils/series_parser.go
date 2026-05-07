package utils

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type ParsedSeries struct {
	Title    string
	Year     int
	Season   int
	Episode  int
	Quality  string
	IsSeries bool
}

// Tokens de release que devem ser removidos ANTES de processar o título
var releaseTokens = []string{
	// Qualidade
	"2160p", "1080p", "720p", "480p", "360p", "4k", "uhd", "hd",
	// Codec
	"x264", "x265", "h264", "h265", "hevc", "avc", "xvid", "x", "h",
	// Source
	"web-dl", "web dl", "webrip", "web rip", "bluray", "blu ray", 
	"brrip", "bdrip", "dvdrip", "hdtv", "hdcam", "cam", "ts",
	// Audio
	"dual", "aac", "ac3", "dts", "5 1", "7 1", "2 0", "atmos",
	// Outros
	"compacto", "extended", "repack", "proper", "unrated", "directors cut",
	// Grupos de release (adicione mais conforme necessário)
	"starckfilmes", "sf", "yify", "rarbg", "etrg", "yts", "amzn", "nf", "netflix",
}

// ParseSeriesFilename tenta detectar se é uma série e extrair informações
func ParseSeriesFilename(path string) ParsedSeries {
	filename := filepath.Base(path)
	dirName := filepath.Base(filepath.Dir(path))
	ext := filepath.Ext(filename)
	filename = strings.TrimSuffix(filename, ext)

	result := ParsedSeries{
		IsSeries: false,
	}

	// Detectar padrões de série
	season, episode, found := extractSeasonEpisode(filename)
	if !found {
		return result
	}

	result.IsSeries = true
	result.Season = season
	result.Episode = episode

	// Extrair qualidade
	result.Quality = extractQuality(filename)

	// Extrair ano (tentar no filename primeiro, depois no diretório)
	result.Year = extractYear(filename)
	if result.Year == 0 {
		result.Year = extractYear(dirName)
	}

	// Extrair título
	// 1. Normalizar separadores PRIMEIRO (. e _ para espaço)
	title := normalizeSeparators(filename)
	
	// 2. Remover season/episode
	title = removeSeasonEpisode(title)
	
	// 3. Remover ano
	title = removeYear(title)
	
	// 4. Remover tokens de release (agora que os pontos viraram espaços)
	title = removeReleaseTokens(title)
	
	// 5. Limpar espaços e capitalizar
	title = cleanSpaces(title)

	result.Title = title

	return result
}

// extractSeasonEpisode detecta vários padrões de season/episode
func extractSeasonEpisode(value string) (season int, episode int, found bool) {
	patterns := []string{
		// S01E01, S1E1
		`[Ss](\d{1,2})[Ee](\d{1,3})`,
		// 1x01, 1x1
		`(\d{1,2})[xX](\d{1,3})`,
		// Season 1 Episode 1
		`[Ss]eason\s*(\d{1,2})\s*[Ee]pisode\s*(\d{1,3})`,
		// 101, 1001 (season 1 episode 1, season 10 episode 1)
		`^(\d{1,2})(\d{2})$`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(value)

		if len(matches) >= 3 {
			season, _ = strconv.Atoi(matches[1])
			episode, _ = strconv.Atoi(matches[2])
			found = true
			return
		}
	}

	return 0, 0, false
}

// removeSeasonEpisode remove padrões de season/episode do título
func removeSeasonEpisode(value string) string {
	patterns := []string{
		`[Ss]\d{1,2}[Ee]\d{1,3}`,
		`\d{1,2}[xX]\d{1,3}`,
		`[Ss]eason\s*\d{1,2}\s*[Ee]pisode\s*\d{1,3}`,
		`\d{3,4}`, // Remove números como 101, 1001
	}

	result := value
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		result = re.ReplaceAllString(result, " ")
	}

	return result
}

// removeReleaseTokens remove tokens de release do título
func removeReleaseTokens(value string) string {
	lower := strings.ToLower(value)

	for _, token := range releaseTokens {
		// Usar word boundary para evitar remover partes de palavras
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(token) + `\b`)
		lower = re.ReplaceAllString(lower, " ")
	}

	// Remover números soltos (1, 2, 3, etc) que sobraram
	re := regexp.MustCompile(`\s+\d+\s+`)
	lower = re.ReplaceAllString(lower, " ")
	
	// Remover números no final
	re = regexp.MustCompile(`\s+\d+$`)
	lower = re.ReplaceAllString(lower, "")

	// Remover letras ASCII soltas (p, x, etc) no meio de espaços
	re = regexp.MustCompile(`\s+[a-z]\s+`)
	lower = re.ReplaceAllString(lower, " ")
	
	// Remover letras ASCII soltas no final
	re = regexp.MustCompile(`\s+[a-z]$`)
	lower = re.ReplaceAllString(lower, "")

	return lower
}

// IsSeriesPath verifica se o caminho parece ser de uma série
func IsSeriesPath(path string) bool {
	filename := filepath.Base(path)
	_, _, found := extractSeasonEpisode(filename)
	return found
}
