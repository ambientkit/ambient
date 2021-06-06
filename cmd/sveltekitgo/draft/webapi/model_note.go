package webapi

import (
	"encoding/json"
	"os"
)

// Note -
type Note struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

var (
	noteStorage = "cmd/sveltekitgo/storage"
)

// LoadNotes -
func LoadNotes(username string) ([]Note, error) {
	b, err := os.ReadFile(noteStorage + "/" + username + ".json")
	if err != nil {
		// If the file can't be read, then ensure unmarshal doesn't crash.
		b = []byte("[]")
	}

	arr := make([]Note, 0)

	err = json.Unmarshal(b, &arr)
	if err != nil {
		return nil, err
	}

	return arr, nil
}

// SaveNotes -
func SaveNotes(username string, notes []Note) error {
	b, err := json.Marshal(notes)
	if err != nil {
		return err
	}

	return os.WriteFile(noteStorage+"/"+username+".json", b, 0644)
}

// CreateNote -
func CreateNote(username string, message string) error {
	notes, err := LoadNotes(username)
	if err != nil {
		return err
	}

	uid, err := generateUUID()
	if err != nil {
		return err
	}

	notes = append(notes, Note{
		ID:      uid,
		Message: message,
	})

	return SaveNotes(username, notes)
}

// DeleteNote -
func DeleteNote(username string, noteID string) error {
	notes, err := LoadNotes(username)
	if err != nil {
		return err
	}

	for i, note := range notes {
		if note.ID == noteID {
			notes = removeNoteItem(notes, i)
			break
		}
	}

	return SaveNotes(username, notes)
}

func removeNoteItem(s []Note, index int) []Note {
	return append(s[:index], s[index+1:]...)
}

// UpdateNote -
func UpdateNote(username string, noteID string, noteMessage string) error {
	notes, err := LoadNotes(username)
	if err != nil {
		return err
	}

	for i, note := range notes {
		if note.ID == noteID {
			notes[i].Message = noteMessage
			break
		}
	}

	return SaveNotes(username, notes)
}
