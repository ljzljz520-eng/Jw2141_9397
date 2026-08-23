package store

import "example.com/xiangzhenfarm/internal/domain"

const notesBucket = "notes"

func (s *Database) SaveNote(note domain.CollaborationNote) error {
	return s.Put(notesBucket, note.ID, note)
}

func (s *Database) GetNote(id string) (domain.CollaborationNote, error) {
	var note domain.CollaborationNote
	err := s.Get(notesBucket, id, &note)
	return note, err
}

func (s *Database) ListNotes() ([]domain.CollaborationNote, error) {
	keys, err := s.Keys(notesBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.CollaborationNote, 0, len(keys))
	for _, key := range keys {
		note, getErr := s.GetNote(key)
		if getErr != nil {
			return nil, getErr
		}
		result = append(result, note)
	}
	return result, nil
}

func (s *Database) NotesForFarmer(id string) ([]domain.CollaborationNote, error) {
	all, err := s.ListNotes()
	if err != nil {
		return nil, err
	}
	result := make([]domain.CollaborationNote, 0)
	for _, note := range all {
		if note.FarmerID == id {
			result = append(result, note)
		}
	}
	return result, nil
}
