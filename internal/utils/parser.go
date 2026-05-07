package utils

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type ParsedMovie struct {
	Title   string
	Year    int
	Quality string
}

var garbageTokens = []string{
	"2160p",
	"1080p",
	"720p",
	"4k",
	"web-dl",
	"webrip",
	"bluray",
	"brrip",
	"x264",
	"x265",
	"h264",
	"h265",
	"hdr",
	"dual",
	"dublado",
	"legendado",
	"aac",
	"5.1",
	"7.1",
	"compacto",
}

func ParseMovieFilename(path string) ParsedMovie {

	filename := filepath.Base(path)

	ext := filepath.Ext(filename)

	filename = strings.TrimSuffix(filename, ext)

	year := extractYear(filename)

	quality := extractQuality(filename)

	title := removeYear(filename)

	title = removeGarbage(title)

	title = normalizeSeparators(title)

	title = cleanSpaces(title)

	return ParsedMovie{
		Title:   title,
		Year:    year,
		Quality: quality,
	}
}

func normalizeSeparators(value string) string {

	value = strings.ReplaceAll(value, ".", " ")
	value = strings.ReplaceAll(value, "_", " ")

	return value
}

func extractYear(value string) int {

	re := regexp.MustCompile(`\b(19\d{2}|20\d{2})\b`)

	match := re.FindStringSubmatch(value)

	if len(match) < 2 {
		return 0
	}

	year, _ := strconv.Atoi(match[1])

	return year
}

func extractQuality(value string) string {
	lower := strings.ToLower(value)

	// Ordem de prioridade: 4K > 2160p > 1080p > 720p
	if strings.Contains(lower, "4k") || strings.Contains(lower, "2160p") {
		return "4K"
	}
	if strings.Contains(lower, "1080p") {
		return "1080p"
	}
	if strings.Contains(lower, "720p") {
		return "720p"
	}

	return "SD"
}

func removeYear(value string) string {

	re := regexp.MustCompile(`\b(19\d{2}|20\d{2})\b`)

	return re.ReplaceAllString(value, "")
}

func removeGarbage(value string) string {

	lower := strings.ToLower(value)

	for _, token := range garbageTokens {

		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(token) + `\b`)

		lower = re.ReplaceAllString(lower, "")
	}

	return lower
}

func cleanSpaces(value string) string {

	re := regexp.MustCompile(`\s+`)

	value = re.ReplaceAllString(value, " ")

	value = strings.TrimSpace(value)

	return titleCase(value)
}

// titleCase capitaliza a primeira letra de cada palavra
func titleCase(value string) string {
	words := strings.Fields(value)
	for i, word := range words {
		if len(word) > 0 {
			// Capitalizar primeira letra, manter o resto como está
			runes := []rune(word)
			runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}
