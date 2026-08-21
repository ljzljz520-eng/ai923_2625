package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"go.etcd.io/bbolt"
	"lawindex/internal/model"
)

var (
	errNotFound   = errors.New("record not found")
	caseBucket    = []byte("cases")
	clientBucket  = []byte("clients")
	reviewBucket  = []byte("reviews")
	archiveBucket = []byte("archives")
)

type Store struct {
	db   *bbolt.DB
	mu   sync.RWMutex
	path string
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("storage path is required")
	}
	db, err := bbolt.Open(path, 0600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, fmt.Errorf("open storage: %w", err)
	}
	s := &Store{db: db, path: path}
	if err = s.init(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) init() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{caseBucket, clientBucket, reviewBucket, archiveBucket} {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}
func (s *Store) Path() string { return s.path }

func encode(v any) ([]byte, error) { return json.Marshal(v) }
func decode(data []byte, v any) error {
	if len(data) == 0 {
		return errNotFound
	}
	return json.Unmarshal(data, v)
}

func (s *Store) PutCase(c model.CaseRecord) error { return s.put(caseBucket, c.ID, c) }
func (s *Store) GetCase(id string) (model.CaseRecord, error) {
	var c model.CaseRecord
	err := s.get(caseBucket, id, &c)
	return c, err
}
func (s *Store) ListCases() ([]model.CaseRecord, error) {
	var out []model.CaseRecord
	err := s.list(caseBucket, func() any { return &model.CaseRecord{} }, func(v any) { out = append(out, *(v.(*model.CaseRecord))) })
	return out, err
}
func (s *Store) PutClient(c model.ClientProfile) error { return s.put(clientBucket, c.ID, c) }
func (s *Store) GetClient(id string) (model.ClientProfile, error) {
	var c model.ClientProfile
	err := s.get(clientBucket, id, &c)
	return c, err
}
func (s *Store) PutReview(r model.ReviewRecord) error { return s.put(reviewBucket, r.ID, r) }
func (s *Store) ListReviews(caseID string) ([]model.ReviewRecord, error) {
	var out []model.ReviewRecord
	err := s.list(reviewBucket, func() any { return &model.ReviewRecord{} }, func(v any) {
		r := *(v.(*model.ReviewRecord))
		if caseID == "" || r.CaseID == caseID {
			out = append(out, r)
		}
	})
	return out, err
}
func (s *Store) PutArchive(a model.ArchiveEntry) error { return s.put(archiveBucket, a.ID, a) }
func (s *Store) GetArchive(id string) (model.ArchiveEntry, error) {
	var a model.ArchiveEntry
	err := s.get(archiveBucket, id, &a)
	return a, err
}

func (s *Store) put(bucket []byte, id string, value any) error {
	if id == "" {
		return errors.New("id is required")
	}
	data, err := encode(value)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucket).Put([]byte(id), data) })
}
func (s *Store) get(bucket []byte, id string, target any) error {
	if id == "" {
		return errors.New("id is required")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error { return decode(tx.Bucket(bucket).Get([]byte(id)), target) })
}
func (s *Store) list(bucket []byte, factory func() any, consume func(any)) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).ForEach(func(_, data []byte) error {
			v := factory()
			if err := decode(data, v); err != nil {
				return err
			}
			consume(v)
			return nil
		})
	})
}
func (s *Store) RemoveCase(id string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(caseBucket).Delete([]byte(id)) })
}
func (s *Store) EnsureDirectory(path string) error {
	if path == "" {
		return errors.New("path is required")
	}
	return os.MkdirAll(path, 0750)
}
