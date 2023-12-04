package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

const (
	baseUrl = "http://localhost:3099"
	createPostfix = "/notes"
	getPostfix = "/notes/%d"	
)

type NoteInfo struct {
	Title string `json:"title"`
	Context string `json:"context"`
	Author string `json:"author"`
	IsPublic bool `json:"is_public"`
}

type Note struct {
	ID int64 `json:"id"`
	Info NoteInfo `json:"info"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func createNote() (Note, error) {
	note := NoteInfo{
		Title: gofakeit.BeerName(),
		Context: gofakeit.IPv4Address(),
		Author: gofakeit.Name(),
		IsPublic: gofakeit.Bool(),
	}

	data, err := json.Marshal(note)
	if err != nil {
		return Note{}, err
	}

	resp, err := http.Post(baseUrl+createPostfix, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return Note{}, err
	}

	defer resp.Body.Close()


	if resp.StatusCode != http.StatusCreated {
		return Note{}, errors.New("failed to create note")
	}

	var createdNote Note
	if err = json.NewDecoder(resp.Body).Decode(&createdNote); err != nil {
		return Note{}, err
	}

	return createdNote, nil
}

func getNote(id int64) (Note, error) {
	url := fmt.Sprintf(baseUrl+getPostfix, id)
	log.Printf("url %s", url)
	
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("failed to get note: ", err)
	}

	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusNotFound {
		return Note{}, errors.New("note not found")
	}

	if resp.StatusCode != http.StatusOK {
		return Note{}, errors.New("failed to get note")
	}

	var note Note
	if err = json.NewDecoder(resp.Body).Decode(&note); err != nil {
		return Note{}, err
	}

	return note, nil
}

func main() {
	note, err := createNote()
	if err != nil {
		log.Fatal("failed to create note: ", err)
	}

	log.Printf("Note created: %+v", note)

	note, err = getNote(note.ID)
	if err != nil {
		log.Fatal("failed to get note: ", err)
	}

	log.Printf("Note info got: %+v", note)
}