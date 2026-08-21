package store

import (
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"os"
	"path/filepath"
	"time"
)

type Health struct {
	Path     string
	Open     bool
	Size     int64
	Modified time.Time
}

func (s *Store) Health() Health {
	info, err := os.Stat(s.path)
	result := Health{Path: s.path, Open: s.db != nil}
	if err == nil {
		result.Size = info.Size()
		result.Modified = info.ModTime()
	}
	return result
}
func (s *Store) Validate() error {
	if s.path == "" {
		return errors.New("storage path is empty")
	}
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{caseBucket, clientBucket, reviewBucket, archiveBucket} {
			if tx.Bucket(bucket) == nil {
				return fmt.Errorf("missing bucket %s", bucket)
			}
		}
		return nil
	})
}
func (s *Store) Backup(target string) error {
	if target == "" {
		return errors.New("backup target is required")
	}
	if filepath.Clean(target) == filepath.Clean(s.path) {
		return errors.New("backup target must differ from source")
	}
	if err := os.MkdirAll(filepath.Dir(target), 0750); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	destination, err := os.Create(target)
	if err != nil {
		return err
	}
	defer destination.Close()
	if err := s.db.View(func(tx *bbolt.Tx) error { _, err := tx.WriteTo(destination); return err }); err != nil {
		return fmt.Errorf("backup database: %w", err)
	}
	return destination.Sync()
}
func (s *Store) RemoveArchive(id string) error {
	if id == "" {
		return errors.New("archive id is required")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(archiveBucket).Delete([]byte(id)) })
}
