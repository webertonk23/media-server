package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"media-server/internal/models"
	"media-server/internal/repositories"
)

type TranscoderService struct {
	mediaFileRepo *repositories.MediaFileRepository
}

func NewTranscoderService() *TranscoderService {
	return &TranscoderService{
		mediaFileRepo: repositories.NewMediaFileRepository(),
	}
}

// StartWorker inicia o loop de transcodificação
func (s *TranscoderService) StartWorker() {
	log.Println("Transcoder Worker iniciado. Aguardando arquivos pendentes...")
	for {
		file, err := s.mediaFileRepo.FindNextPending()
		if err != nil {
			// Sem arquivos pendentes, espera um pouco
			time.Sleep(10 * time.Second)
			continue
		}

		log.Printf("Novo arquivo para transcodificar encontrado: %s\n", file.Path)
		err = s.ProcessFile(file)
		if err != nil {
			log.Printf("Erro ao processar arquivo %d (%s): %v\n", file.ID, file.Path, err)
			s.mediaFileRepo.UpdateStatus(file.ID, models.FileStatusError, err.Error())
		}
	}
}

// ProcessFile executa o transcode de um arquivo individual
func (s *TranscoderService) ProcessFile(file *models.MediaFile) error {
	// Marcar como processando imediatamente
	err := s.mediaFileRepo.UpdateStatus(file.ID, models.FileStatusProcessing, "")
	if err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(file.Path))
	dir := filepath.Dir(file.Path)
	
	// Se já for MP4, podemos apenas marcar como concluído ou verificar o codec
	// Por simplicidade agora, vamos transcodificar tudo para garantir h264/aac
	
	// Configurações via variáveis de ambiente ou defaults
	codec := os.Getenv("TRANSCODE_CODEC")
	if codec == "" {
		codec = "libx264" // Default software
	}

	preset := os.Getenv("TRANSCODE_PRESET")
	if preset == "" {
		preset = "fast"
	}

	tempOutputPath := filepath.Join(dir, fmt.Sprintf("transcoding_%s.mp4", file.ULID))

	log.Printf("Iniciando FFmpeg para %s (Codec: %s)\n", file.Path, codec)

	// Montar argumentos do FFmpeg de forma dinâmica
	args := []string{
		"-i", file.Path,
		"-c:v", codec,
	}

	// Presets só existem no libx264, codecs de hardware usam outros parâmetros
	if strings.Contains(codec, "libx264") {
		args = append(args, "-preset", preset, "-crf", "23")
	}

	args = append(args,
		"-pix_fmt", "yuv420p",
		"-c:a", "aac",
		"-b:a", "128k",
		"-y",
		tempOutputPath,
	)

	cmd := exec.Command("ffmpeg", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v, output: %s", err, string(output))
	}

	// Verificar se o arquivo foi criado e tem tamanho
	stat, err := os.Stat(tempOutputPath)
	if err != nil || stat.Size() == 0 {
		return fmt.Errorf("arquivo de saída não gerado ou vazio")
	}

	// Sucesso!
	originalPath := file.Path
	
	// Novo nome final: mesmo nome mas com .mp4
	finalPath := filepath.Join(dir, strings.TrimSuffix(filepath.Base(originalPath), ext)+".mp4")

	// Se o finalPath for igual ao originalPath (ex: era um mp4 com codec ruim)
	// Precisamos de um passo intermediário
	if finalPath == originalPath {
		// Já estamos usando tempOutputPath, então primeiro removemos o original
		err = os.Remove(originalPath)
		if err != nil {
			return fmt.Errorf("falha ao remover original: %v", err)
		}
	}

	// Mover temp para final
	err = os.Rename(tempOutputPath, finalPath)
	if err != nil {
		return fmt.Errorf("falha ao renomear arquivo final: %v", err)
	}

	// Se não era o mesmo nome, precisamos remover o original agora
	if finalPath != originalPath {
		err = os.Remove(originalPath)
		if err != nil {
			log.Printf("Aviso: não foi possível remover original %s: %v\n", originalPath, err)
		}
	}

	// Atualizar banco
	file.Path = finalPath
	file.OriginalPath = originalPath
	file.Status = models.FileStatusCompleted
	file.Size = stat.Size()
	
	err = s.mediaFileRepo.Update(file)
	if err != nil {
		return err
	}

	log.Printf("Arquivo transcodificado com sucesso: %s\n", finalPath)
	return nil
}
