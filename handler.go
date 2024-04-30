package main

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	Store      *MongoStore
	InMemoryDb *InMemoryDb
}

func newHandler(store *MongoStore, inMemoryDb *InMemoryDb) *Handler {
	return &Handler{
		Store:      store,
		InMemoryDb: inMemoryDb,
	}
}

func (h *Handler) HandleGetRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	start, err := timeFromString(req.StartDate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	end, err := timeFromString(req.EndDate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	records, err := h.Store.GetRecords(start, end, req.MinCount, req.MaxCount)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(*records) == 0 {
		respondWithError(w, http.StatusNotFound, "No records found")
		return
	}

	var respRecords []RecordDTO

	for _, record := range *records {
		respRecords = append(respRecords, *ResponseFromRecord(record))
	}

	resp := Response{
		Count:   len(respRecords),
		Msg:     "success",
		Records: respRecords,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (h *Handler) HandlePostPayload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	type Req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	var req Req

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := h.InMemoryDb.Insert(req.Key, req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (h *Handler) HandleGetPayloadById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	queryParams := r.URL.Query()
	key := queryParams.Get("key")

	resp, err := h.InMemoryDb.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}
