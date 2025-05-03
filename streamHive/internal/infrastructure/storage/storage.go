package storage

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Storage для хранит видеофайлы локально и удаляет их по TTL
type Storage struct {
	path      string
	retention time.Duration
	mu        sync.Mutex
}

func New(path string, retention string) *Storage {
	dur, err := time.ParseDuration(retention)
	if err != nil {
		log.Fatalf("invalid retention duration: %v", err)
	}
	return &Storage{
		path:      path,
		retention: dur,
	}
}

func (s *Storage) Save(jobID string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	filename := filepath.Join(s.path, jobID+".mp4")
	return os.WriteFile(filename, data, 0644)
}

func (s *Storage) StartCleanupWorker() {
	go func() {
		for {
			_ = filepath.Walk(s.path, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				if time.Since(info.ModTime()) > s.retention {
					log.Printf("removing expired file: %s", path)
					_ = os.Remove(path)
				}
				return nil
			})
			time.Sleep(1 * time.Hour)
		}
	}()
}
