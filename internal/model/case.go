package model

import "time"

const (
	EntityCaseRecord    = "CaseRecord"
	EntityClientProfile = "ClientProfile"
	EntityReviewRecord  = "ReviewRecord"
	EntityArchiveEntry  = "ArchiveEntry"
)

type CaseRecord struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	ClientID  string    `json:"client_id"`
	Category  string    `json:"category"`
	Summary   string    `json:"summary"`
	Status    string    `json:"status"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ClientProfile struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Notes   string `json:"notes"`
}

type ReviewRecord struct {
	ID        string    `json:"id"`
	CaseID    string    `json:"case_id"`
	Reviewer  string    `json:"reviewer"`
	Decision  string    `json:"decision"`
	Comment   string    `json:"comment"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

type ArchiveEntry struct {
	ID         string    `json:"id"`
	CaseID     string    `json:"case_id"`
	Reason     string    `json:"reason"`
	ArchivedBy string    `json:"archived_by"`
	ArchivedAt time.Time `json:"archived_at"`
}

type SearchFilter struct{ Query, Category, Status string }

func (c CaseRecord) IsOpen() bool         { return c.Status != "archived" }
func (c CaseRecord) IsReviewed() bool     { return c.Status == "approved" || c.Status == "rejected" }
func (c CaseRecord) CanArchive() bool     { return c.IsReviewed() && c.Status != "archived" }
func (c *CaseRecord) Touch(now time.Time) { c.UpdatedAt = now }
