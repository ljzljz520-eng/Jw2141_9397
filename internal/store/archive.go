package store

import "example.com/xiangzhenfarm/internal/domain"

const archivesBucket = "archives"

func (s *Database) SaveArchive(entry domain.ArchiveEntry) error {
	return s.Put(archivesBucket, entry.ID, entry)
}

func (s *Database) GetArchive(id string) (domain.ArchiveEntry, error) {
	var entry domain.ArchiveEntry
	err := s.Get(archivesBucket, id, &entry)
	return entry, err
}

func (s *Database) ListArchives() ([]domain.ArchiveEntry, error) {
	keys, err := s.Keys(archivesBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ArchiveEntry, 0, len(keys))
	for _, key := range keys {
		entry, getErr := s.GetArchive(key)
		if getErr != nil {
			return nil, getErr
		}
		result = append(result, entry)
	}
	return result, nil
}

func (s *Database) ArchivesForFarmer(id string) ([]domain.ArchiveEntry, error) {
	all, err := s.ListArchives()
	if err != nil {
		return nil, err
	}
	result := make([]domain.ArchiveEntry, 0)
	for _, entry := range all {
		if entry.FarmerID == id {
			result = append(result, entry)
		}
	}
	return result, nil
}
