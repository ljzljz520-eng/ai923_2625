package review

import (
	"lawindex/internal/model"
	"sort"
)

type QueueItem struct {
	Case    model.CaseRecord
	Reviews int
}

func (s *Service) Queue() ([]QueueItem, error) {
	cases, err := s.storage.ListCases()
	if err != nil {
		return nil, err
	}
	result := []QueueItem{}
	for _, item := range cases {
		if item.Status != "submitted" {
			continue
		}
		reviews, reviewErr := s.storage.ListReviews(item.ID)
		if reviewErr != nil {
			return nil, reviewErr
		}
		result = append(result, QueueItem{Case: item, Reviews: len(reviews)})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Case.CreatedAt.Before(result[j].Case.CreatedAt) })
	return result, nil
}
func (s *Service) PendingCount() int {
	queue, err := s.Queue()
	if err != nil {
		return 0
	}
	return len(queue)
}
func SortQueue(items []QueueItem, newestFirst bool) []QueueItem {
	result := append([]QueueItem(nil), items...)
	sort.Slice(result, func(i, j int) bool {
		if newestFirst {
			return result[i].Case.CreatedAt.After(result[j].Case.CreatedAt)
		}
		return result[i].Case.CreatedAt.Before(result[j].Case.CreatedAt)
	})
	return result
}
func QueueIDs(items []QueueItem) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.Case.ID)
	}
	return ids
}
