package store

import (
	"encoding/json"
	"errors"
	"lawindex/internal/model"
)

type Snapshot struct {
	Cases    []model.CaseRecord
	Clients  []model.ClientProfile
	Reviews  []model.ReviewRecord
	Archives []model.ArchiveEntry
}

func (s *Store) Snapshot() (Snapshot, error) {
	cases, err := s.ListCases()
	if err != nil {
		return Snapshot{}, err
	}
	snap := Snapshot{Cases: cases, Clients: []model.ClientProfile{}, Reviews: []model.ReviewRecord{}, Archives: []model.ArchiveEntry{}}
	for _, c := range cases {
		if client, e := s.GetClient(c.ClientID); e == nil {
			snap.Clients = append(snap.Clients, client)
		}
		rs, e := s.ListReviews(c.ID)
		if e != nil {
			return Snapshot{}, e
		}
		snap.Reviews = append(snap.Reviews, rs...)
	}
	return snap, nil
}
func (s *Store) SnapshotJSON() ([]byte, error) {
	snap, err := s.Snapshot()
	if err != nil {
		return nil, err
	}
	return json.Marshal(snap)
}
func DecodeSnapshot(data []byte) (Snapshot, error) {
	if len(data) == 0 {
		return Snapshot{}, errors.New("snapshot is empty")
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}
func (s *Store) RestoreSnapshot(snap Snapshot) error {
	for _, c := range snap.Cases {
		if err := s.PutCase(c); err != nil {
			return err
		}
	}
	for _, c := range snap.Clients {
		if err := s.PutClient(c); err != nil {
			return err
		}
	}
	for _, r := range snap.Reviews {
		if err := s.PutReview(r); err != nil {
			return err
		}
	}
	for _, a := range snap.Archives {
		if err := s.PutArchive(a); err != nil {
			return err
		}
	}
	return nil
}
