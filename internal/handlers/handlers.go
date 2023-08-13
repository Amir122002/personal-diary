package handlers

import (
	"diary/pkg/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var diary models.Diary
var maxID int

func ReadAll(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(diary.Notes)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(response)
	if err != nil {
		return
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	var notesData map[string]interface{}
	err = json.Unmarshal(body, &notesData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	text, ok := notesData["text"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newNote := createNewNote(text)

	diary.Notes = append(diary.Notes, newNote)
	saveNotes()

	w.WriteHeader(http.StatusCreated)
}

func createNewNote(text string) models.Note {
	maxID++
	date := time.Now().Format("2006-01-02 15:04:05")
	return models.Note{
		ID:   maxID,
		Text: text,
		Date: date,
	}
}

func saveNotes() error {
	data, err := json.Marshal(diary.Notes)
	if err != nil {
		return err
	}

	file, err := os.OpenFile("notebooks.json", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func LoadNotes() {
	data, err := os.ReadFile("./internal/notebooks/notes.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &diary.Notes)
	if err != nil {
		fmt.Println("Error during JSON unmarshalling:", err)
		return
	}
	maxID = 0
	for _, note := range diary.Notes {
		if note.ID > maxID {
			maxID = note.ID
		}
	}
}

func HandleNote(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/notes/"):]
	id, _ := strconv.Atoi(idStr)

	for i, note := range diary.Notes {
		if note.ID == id {
			switch {

			case r.Method == http.MethodGet:
				response, _ := json.Marshal(note)
				w.Header().Set("Content-Type", "application/json")
				_, err := w.Write(response)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				return
			case r.Method == http.MethodPut:
				body, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body.Close()
				var notesData map[string]interface{}
				err = json.Unmarshal(body, &notesData)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				text, ok := notesData["text"].(string)
				if !ok {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				diary.Notes[i].Text = text
				err = saveNotes()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				return
			case r.Method == http.MethodDelete:
				diary.Notes = append(diary.Notes[:i], diary.Notes[i+1:]...)
				err := saveNotes()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	return
}
