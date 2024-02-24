package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type OpenSourceProject struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	OpenIssues []string  `json:"open_issues"`
	OpenPRs    []string  `json:"open_prs"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateOpenSourceProjectReq struct {
	Name       string   `json:"name"`
	OpenIssues []string `json:"open_issues"`
	OpenPRs    []string `json:"open_prs"`
}

type projectHandlers struct {
	sync.Mutex
	db map[string]OpenSourceProject
}

func (h *projectHandlers) projects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getAll(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}

func (h *projectHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type application-json, but got %s", ct)))
		return
	}

	var body CreateOpenSourceProjectReq
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	openSourceProject := OpenSourceProject{
		ID:         fmt.Sprint(len(h.db) + 1),
		Name:       body.Name,
		OpenIssues: body.OpenIssues,
		OpenPRs:    body.OpenPRs,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	h.Lock()
	h.db[openSourceProject.ID] = openSourceProject
	h.Unlock()
}

func (h *projectHandlers) getAll(w http.ResponseWriter, r *http.Request) {
	projects := make([]OpenSourceProject, len(h.db))

	h.Lock()
	i := 0
	for _, project := range h.db {
		projects[i] = project
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(projects)
	fmt.Printf("json bytes: %v", jsonBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(jsonBytes)
}

func newProjectHandlers() *projectHandlers {
	return &projectHandlers{
		db: map[string]OpenSourceProject{
			"1": {
				ID:         "1",
				Name:       "Project 1",
				OpenIssues: []string{"1", "2"},
				OpenPRs:    []string{"1", "2"},
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			"2": {
				ID:         "2",
				Name:       "Project 2",
				OpenIssues: []string{"1", "2"},
				OpenPRs:    []string{"1", "2"},
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			"3": {
				ID:         "3",
				Name:       "Project 3",
				OpenIssues: []string{"1", "2"},
				OpenPRs:    []string{"1", "2"},
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
	}
}
func main() {
	fmt.Println("Start server")

	openSourceHandlers := newProjectHandlers()
	http.HandleFunc("/opensource/projects", openSourceHandlers.projects)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
