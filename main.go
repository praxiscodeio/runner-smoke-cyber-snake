package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var DefaultRules = map[string]string{
	".jpg":  "Images",
	".jpeg": "Images",
	".png":  "Images",
	".pdf":  "Documents",
	".doc":  "Documents",
	".docx": "Documents",
	".txt":  "Documents",
	".mp3":  "Music",
	".wav":  "Music",
	".mp4":  "Video",
	".avi":  "Video",
	".zip":  "Archives",
	".rar":  "Archives",
}

type FileOrganizer struct {
	sourceDir      string
	rulesMap       map[string]string
	processedFiles int
	logFile        *os.File
}

func NewFileOrganizer(sourceDir string) (*FileOrganizer, error) {
	if sourceDir == "" {
		return nil, fmt.Errorf("путь к директории не может быть пустым")
	}
	// «inbox/» и «./inbox» — та же директория; Organize сравнивает с ней filepath.Dir
	sourceDir = filepath.Clean(sourceDir)

	info, err := os.Stat(sourceDir)
	if err != nil {
		return nil, fmt.Errorf("директория не найдена: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("указанный путь не является директорией: %s", sourceDir)
	}

	return &FileOrganizer{
		sourceDir: sourceDir,
		rulesMap:  DefaultRules,
	}, nil
}

func main() {
	fo, err := NewFileOrganizer(".")
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	fmt.Printf("FileOrganizer создан для директории: %s\n", fo.sourceDir)
}
