package models
import "gorm.io/gorm"
type Settings struct {
	gorm.Model
	MoviePath    string `json:"movie_path"`
	SeriesPath   string `json:"series_path"`
	ScanInterval int    `json:"scan_interval"`
}
