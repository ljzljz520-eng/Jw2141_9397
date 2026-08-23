package store

import "example.com/xiangzhenfarm/internal/domain"

const farmersBucket = "farmers"

func (s *Database) SaveFarmer(record domain.FarmerRecord) error {
	return s.Put(farmersBucket, record.ID, record)
}

func (s *Database) GetFarmer(id string) (domain.FarmerRecord, error) {
	var record domain.FarmerRecord
	err := s.Get(farmersBucket, id, &record)
	return record, err
}

func (s *Database) ListFarmers() ([]domain.FarmerRecord, error) {
	keys, err := s.Keys(farmersBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.FarmerRecord, 0, len(keys))
	for _, key := range keys {
		record, getErr := s.GetFarmer(key)
		if getErr != nil {
			return nil, getErr
		}
		result = append(result, record)
	}
	return result, nil
}

func (s *Database) DeleteFarmer(id string) error {
	return s.Delete(farmersBucket, id)
}

func (s *Database) FarmerCount() (int, error) {
	return s.Count(farmersBucket)
}
