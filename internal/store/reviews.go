package store

import "example.com/xiangzhenfarm/internal/domain"

const (
	visitsBucket  = "visits"
	reviewsBucket = "reviews"
)

func (s *Database) SaveVisitReport(report domain.VisitReport) error {
	return s.Put(visitsBucket, report.ID, report)
}

func (s *Database) GetVisitReport(id string) (domain.VisitReport, error) {
	var report domain.VisitReport
	err := s.Get(visitsBucket, id, &report)
	return report, err
}

func (s *Database) SaveReview(review domain.ReviewCase) error {
	return s.Put(reviewsBucket, review.ID, review)
}

func (s *Database) GetReview(id string) (domain.ReviewCase, error) {
	var review domain.ReviewCase
	err := s.Get(reviewsBucket, id, &review)
	return review, err
}

func (s *Database) ListReviews() ([]domain.ReviewCase, error) {
	keys, err := s.Keys(reviewsBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ReviewCase, 0, len(keys))
	for _, key := range keys {
		review, getErr := s.GetReview(key)
		if getErr != nil {
			return nil, getErr
		}
		result = append(result, review)
	}
	return result, nil
}
