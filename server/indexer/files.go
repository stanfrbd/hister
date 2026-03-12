package indexer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/rs/zerolog/log"
)

const maxFileSize = 1024 * 1024 // 1MB

func WatchDirectories(dirs []string, interval time.Duration) {
	indexAll(dirs)
	ticker := time.NewTicker(interval)
	for range ticker.C {
		indexAll(dirs)
	}
}

func indexAll(dirs []string) {
	for _, dir := range dirs {
		dir = expandHome(dir)
		if err := indexDirectory(dir); err != nil {
			log.Error().Err(err).Str("directory", dir).Msg("Failed to index directory")
		}
	}
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func indexDirectory(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("cannot access directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", dir)
	}

	indexed := 0
	skipped := 0

	log.Debug().Str("directory", dir).Msg("Indexing directory")

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Warn().Err(err).Str("path", path).Msg("Error accessing path")
			return nil
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") && path != dir {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}
		if err := indexFile(path); err != nil {
			log.Debug().Err(err).Str("path", path).Msg("Skipping file")
			skipped++
		} else {
			indexed++
		}
		return nil
	})

	log.Info().Str("directory", dir).Int("indexed", indexed).Int("skipped", skipped).Msg("Directory indexing complete")
	return err
}

func indexFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		return fmt.Errorf("empty file")
	}
	if info.Size() > maxFileSize {
		return fmt.Errorf("file too large (%d bytes)", info.Size())
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	fileURL := "file://" + absPath

	// Skip if already indexed with the same modification time
	existing := GetByURL(fileURL)
	if existing != nil && existing.Added == info.ModTime().Unix() {
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}
	if !utf8.Valid(content) {
		return fmt.Errorf("binary file")
	}

	doc := &Document{
		URL:   fileURL,
		Text:  string(content),
		Added: info.ModTime().Unix(),
	}

	return i.AddDocument(doc)
}
