package archive

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"lawindex/internal/model"
	"lawindex/internal/store"
)

type Service struct {
	storage *store.Store
	now     func() time.Time
}

func New(storage *store.Store) *Service { return &Service{storage: storage, now: time.Now} }
func (s *Service) Archive(caseID, actor, reason string) (model.ArchiveEntry, error) {
	c, err := s.storage.GetCase(caseID)
	if err != nil {
		return model.ArchiveEntry{}, err
	}
	if !c.CanArchive() {
		return model.ArchiveEntry{}, errors.New("case must be reviewed before archive")
	}
	c.Status = "archived"
	c.UpdatedAt = s.now()
	if err := s.storage.PutCase(c); err != nil {
		return model.ArchiveEntry{}, err
	}
	a := model.ArchiveEntry{ID: fmt.Sprintf("archive-%d", s.now().UnixNano()), CaseID: caseID, Reason: strings.TrimSpace(reason), ArchivedBy: strings.TrimSpace(actor), ArchivedAt: s.now()}
	return a, s.storage.PutArchive(a)
}
func (s *Service) Restore(caseID string) error {
	c, err := s.storage.GetCase(caseID)
	if err != nil {
		return err
	}
	if c.Status != "archived" {
		return errors.New("case is not archived")
	}
	c.Status = "approved"
	c.UpdatedAt = s.now()
	return s.storage.PutCase(c)
}
func (s *Service) ArchiveDetails(id string) (model.ArchiveEntry, error) {
	return s.storage.GetArchive(id)
}
