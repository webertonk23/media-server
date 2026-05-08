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

// ParseSeriesFilename tenta detectar se é uma série e extrair informações
func extractSeasonEpisodeWithLoc(value string) (season int, episode int, loc []int, found bool) {
	// 1. Limpar tags entre colchetes no início (comum em animes)
	// Ex: [Erai-raws] Made in Abyss - 01 -> Made in Abyss - 01
	reBrackets := regexp.MustCompile(`^\[[^\]]+\]\s*`)
	tagLoc := reBrackets.FindStringIndex(value)
	searchTarget := value
	offset := 0
	if tagLoc != nil {
		searchTarget = value[tagLoc[1]:]
		offset = tagLoc[1]
	}

	patterns := []string{
		// S01E01, S1E1
		`(?i)[.\s\(\[]*[Ss](\d{1,2})[Ee](\d{1,3})`,
		// 1x01, 1x1
		`(?i)[.\s\(\[]*(\d{1,2})[xX](\d{1,3})`,
		// - 01 (Hyphenated anime style)
		`(?i)[.\s\(\[]*-\s*(\d{1,4})`,
		// E1066 (Anime style)
		`(?i)[.\s\(\[]*[Ee](\d{2,4})`,
		// Season 1 Episode 1
		`(?i)[.\s\(\[]*[Ss]eason\s*(\d{1,2})\s*[Ee]pisode\s*(\d{1,3})`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matchLoc := re.FindStringIndex(searchTarget)
		if matchLoc != nil {
			matches := re.FindStringSubmatch(searchTarget[matchLoc[0]:matchLoc[1]])
			if len(matches) >= 3 {
				season, _ = strconv.Atoi(matches[1])
				episode, _ = strconv.Atoi(matches[2])
				found = true
				loc = []int{matchLoc[0] + offset, matchLoc[1] + offset}
				return
			} else if len(matches) == 2 {
				// Caso E1066 ou - 01
				season = 1 // Default
				episode, _ = strconv.Atoi(matches[1])
				found = true
				loc = []int{matchLoc[0] + offset, matchLoc[1] + offset}
				return
			}
		}
	}

	return 0, 0, nil, false
}

func ParseSeriesFilename(path string) ParsedSeries {
	filename := filepath.Base(path)
	dirName := filepath.Base(filepath.Dir(path))
	ext := filepath.Ext(filename)
	cleanName := strings.TrimSuffix(filename, ext)

	// 1. Detectar padrões de série e pegar a localização
	season, episode, loc, found := extractSeasonEpisodeWithLoc(cleanName)
	if !found {
		return ParsedSeries{IsSeries: false}
	}

	result := ParsedSeries{
		IsSeries: true,
		Season:   season,
		Episode:  episode,
	}

	// 2. Título é tudo antes do marcador de S/E
	titleRaw := cleanName[:loc[0]]

	// 3. Remover tags entre colchetes do título (ex: [Erai-raws])
	reBrackets := regexp.MustCompile(`\[[^\]]+\]`)
	titleRaw = reBrackets.ReplaceAllString(titleRaw, "")

	// 4. Extrair qualidade e ano
	result.Quality = extractQuality(cleanName)
	result.Year = extractYear(cleanName)
	if result.Year == 0 {
		result.Year = extractYear(dirName)
	}

	// 5. Se o título cru contém o ano, remover o ano e o que vem depois dele
	reYear := regexp.MustCompile(`(?i)[.\s\(\[]+(19\d{2}|20\d{2})[.\s\)\]]*`)
	yearLoc := reYear.FindStringIndex(titleRaw)
	if yearLoc != nil {
		titleRaw = titleRaw[:yearLoc[0]]
	}

	// 6. Limpeza final
	title := normalizeSeparators(titleRaw)
	title = cleanSpaces(title)
	result.Title = title

	return result
}

// IsSeriesPath verifica se o caminho parece ser de uma série
func IsSeriesPath(path string) bool {
	filename := filepath.Base(path)
	_, _, _, found := extractSeasonEpisodeWithLoc(filename)
	return found
}
