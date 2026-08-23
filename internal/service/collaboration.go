package service

import (
	"fmt"
	"sort"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/validate"
)

func (r *Registry) AddCollaborationNote(note domain.CollaborationNote) (domain.CollaborationNote, error) {
	note.Author = strings.TrimSpace(note.Author)
	note.Body = strings.TrimSpace(note.Body)
	note.Visibility = strings.TrimSpace(strings.ToLower(note.Visibility))
	if note.ID == "" {
		note.ID = domain.KeyFor("note", note.FarmerID+"|"+note.Author+"|"+note.Body)
	}
	if note.Visibility == "" {
		note.Visibility = "team"
	}
	if note.CreatedAt == "" {
		note.CreatedAt = "2026-01-01"
	}
	if validation := validate.CollaborationNote(note); validate.HasErrors(validation) {
		return domain.CollaborationNote{}, fmt.Errorf("validate collaboration note: %s", strings.Join(validation, "; "))
	}
	if _, err := r.GetFarmer(note.FarmerID); err != nil {
		return domain.CollaborationNote{}, err
	}
	if err := r.store.SaveNote(note); err != nil {
		return domain.CollaborationNote{}, err
	}
	if _, err := r.RecordAudit("CollaborationNote", note.ID, domain.AuditNoteAdded, note.Author, note.CreatedAt, note.Body); err != nil {
		return domain.CollaborationNote{}, err
	}
	return note, nil
}

func (r *Registry) FarmerNotes(id string) ([]domain.CollaborationNote, error) {
	notes, err := r.store.NotesForFarmer(id)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(notes, func(left, right int) bool {
		return notes[left].ID < notes[right].ID
	})
	return notes, nil
}

func (r *Registry) CollaborationDigest(id string) (string, error) {
	notes, err := r.FarmerNotes(id)
	if err != nil {
		return "", err
	}
	lines := make([]string, 0, len(notes))
	for _, note := range notes {
		if note.Visibility == "private" {
			continue
		}
		lines = append(lines, note.Author+": "+note.Body)
	}
	return strings.Join(lines, "\n"), nil
}

func (r *Registry) NoteCount(id string) (int, error) {
	notes, err := r.FarmerNotes(id)
	return len(notes), err
}
