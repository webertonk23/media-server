package metadata

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"media-server/internal/config"
)

type TMDBSearchResponse struct {
	Results []TMDBMovie `json:"results"`
}

type TMDBMovie struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	OriginalTitle string `json:"original_title"`
	Overview      string `json:"overview"`
	PosterPath    string `json:"poster_path"`
	Backdrop      string `json:"backdrop_path"`
	ReleaseDate   string `json:"release_date"`
	Year          int    // Calculado a partir de ReleaseDate
}

type TMDBSeries struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	Backdrop     string `json:"backdrop_path"`
	FirstAirDate string `json:"first_air_date"`
	Year         int
}

type TMDBSeriesSearchResponse struct {
	Results []TMDBSeries `json:"results"`
}

func SearchMovie(title string, year int) (*TMDBMovie, error) {
	baseURL := "https://api.themoviedb.org/3/search/movie"
	params := url.Values{}
	params.Add("api_key", config.AppConfig.TMDBApiKey)
	params.Add("query", title)
	params.Add("language", "pt-BR")

	if year > 0 {
		params.Add("year", fmt.Sprintf("%d", year))
	}

	fullURL := baseURL + "?" + params.Encode()
	log.Println(fullURL)

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TMDBSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if len(result.Results) == 0 {
		return nil, nil
	}

	movie := &result.Results[0]
	if len(movie.ReleaseDate) >= 4 {
		fmt.Sscanf(movie.ReleaseDate[:4], "%d", &movie.Year)
	}

	return movie, nil
}

func SearchSeries(title string, year int) (*TMDBSeries, error) {
	baseURL := "https://api.themoviedb.org/3/search/tv"
	params := url.Values{}
	params.Add("api_key", config.AppConfig.TMDBApiKey)
	params.Add("query", title)
	params.Add("language", "pt-BR")

	fullURL := baseURL + "?" + params.Encode()
	log.Println(fullURL)

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TMDBSeriesSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if len(result.Results) == 0 {
		return nil, nil
	}

	series := &result.Results[0]
	if len(series.FirstAirDate) >= 4 {
		fmt.Sscanf(series.FirstAirDate[:4], "%d", &series.Year)
	}

	return series, nil
}
