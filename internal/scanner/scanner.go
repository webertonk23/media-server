package scanner

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ScannedFile representa um arquivo descoberto pelo scanner
// O scanner NÃO conhece TMDB, Movie ou Series
type ScannedFile struct {
	Filename    string
	Path        string
	Size        int64
	ModifiedAt  time.Time
	Fingerprint string
	Extension   string
}

// SupportedExtensions define as extensões de vídeo suportadas
var SupportedExtensions = []string{
	".mp4",
	".mkv",
	".avi",
	".mov",
	".wmv",
	".flv",
	".webm",
	".m4v",
}

// ScanDirectory percorre um diretório e retorna todos os arquivos de vídeo encontrados
func ScanDirectory(rootPath string) ([]ScannedFile, error) {
	var files []ScannedFile

	err := filepath.Walk(
		rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			ext := strings.ToLower(filepath.Ext(path))

			if !isSupportedExtension(ext) {
				return nil
			}

			filename := strings.TrimSuffix(info.Name(), ext)

			fingerprint := generateFingerprint(
				path,
				info.Size(),
				info.ModTime().String(),
			)

			files = append(files, ScannedFile{
				Filename:    filename,
				Path:        path,
				Size:        info.Size(),
				ModifiedAt:  info.ModTime(),
				Fingerprint: fingerprint,
				Extension:   ext,
			})

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return files, nil
}

func isSupportedExtension(ext string) bool {
	for _, supported := range SupportedExtensions {
		if ext == supported {
			return true
		}
	}
	return false
}

func generateFingerprint(
	path string,
	size int64,
	modified string,
) string {
	hash := md5.Sum([]byte(
		path + modified,
	))
	return hex.EncodeToString(hash[:])
}
