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
func extractSeasonEpisodeWithLoc(value string) (season int, episode int, loc []int, found bool) {
	reBrackets := regexp.MustCompile(`^\[[^\]]+\]\s*`)
	tagLoc := reBrackets.FindStringIndex(value)
	searchTarget := value
	offset := 0
	if tagLoc != nil {
		searchTarget = value[tagLoc[1]:]
		offset = tagLoc[1]
	}
	patterns := []string{
		`(?i)[.\s\(\[]*[Ss](\d{1,2})[Ee](\d{1,3})`,
		`(?i)[.\s\(\[]*(\d{1,2})[xX](\d{1,3})`,
		`(?i)[.\s\(\[]*-\s*(\d{1,4})`,
		`(?i)[.\s\(\[]*[Ee](\d{2,4})`,
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
				season = 1 
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
	season, episode, loc, found := extractSeasonEpisodeWithLoc(cleanName)
	if !found {
		return ParsedSeries{IsSeries: false}
	}
	result := ParsedSeries{
		IsSeries: true,
		Season:   season,
		Episode:  episode,
	}
	titleRaw := cleanName[:loc[0]]
	reBrackets := regexp.MustCompile(`\[[^\]]+\]`)
	titleRaw = reBrackets.ReplaceAllString(titleRaw, "")
	result.Quality = extractQuality(cleanName)
	result.Year = extractYear(cleanName)
	if result.Year == 0 {
		result.Year = extractYear(dirName)
	}
	reYear := regexp.MustCompile(`(?i)[.\s\(\[]+(19\d{2}|20\d{2})[.\s\)\]]*`)
	yearLoc := reYear.FindStringIndex(titleRaw)
	if yearLoc != nil {
		titleRaw = titleRaw[:yearLoc[0]]
	}
	title := normalizeSeparators(titleRaw)
	title = cleanSpaces(title)
	result.Title = title
	return result
}
func IsSeriesPath(path string) bool {
	filename := filepath.Base(path)
	_, _, _, found := extractSeasonEpisodeWithLoc(filename)
	return found
}
