package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	baseUrl = "localhost:3099"
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

type SyncMap struct {
	elems map[int64]*Note
	m sync.RWMutex
}

var notes = &SyncMap{
	elems: make(map[int64]*Note),
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	info := &NoteInfo{}

	if err := json.NewDecoder(r.Body).Decode(info); err != nil {
		http.Error(w, "failed to decode note data", http.StatusBadRequest)
		return
	}

	rand.Seed(time.Now().UnixNano())
	now := time.Now()

	note := &Note{
		ID: rand.Int63(),
		Info: *info,
		CreatedAt: now,
		UpdatedAt: now,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "failed to encode note data", http.StatusInternalServerError)
		return
	}

	notes.m.Lock()
	defer notes.m.Unlock()

	notes.elems[note.ID] = note
}

func getNoteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("hit get")	
	
	noteID := chi.URLParam(r, "id")
	id, err := parseNoteID(noteID)

	log.Printf("result: %+v", notes.elems)
	log.Printf("got id %d: ", id)

	if err != nil {
		http.Error(w, "invalid note id", http.StatusBadRequest)
	}

	notes.m.RLock()
	defer notes.m.RUnlock()

	note, ok := notes.elems[id]
	if !ok {
		http.Error(w, "note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "failed to encode note data", http.StatusInternalServerError)
		return
	}
}

func parseNoteID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, nil
	}

	return id, nil
}
 
func main() {
	r := chi.NewRouter()

	r.Post(createPostfix, createNoteHandler)
	r.Get(getPostfix, getNoteHandler)

	fmt.Println("Starting http server...")
	err := http.ListenAndServe(baseUrl, r)
	if err != nil {
		log.Fatal(err)
	}
}