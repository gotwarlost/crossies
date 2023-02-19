package api

import (
	"encoding/json"
	"net/http"

	"github.com/gotwarlost/crossies/internal/anagrams"
	"github.com/gotwarlost/crossies/internal/findwords"
	"github.com/gotwarlost/crossies/internal/inputerror"
	"github.com/gotwarlost/crossies/internal/synonyms"
)

type Handler struct {
	h http.Handler
}

func New() (*Handler, error) {
	mux := http.NewServeMux()
	ret := &Handler{}
	mux.Handle("/v1/synonyms", http.HandlerFunc(ret.synonyms))
	mux.Handle("/v1/matching-words", http.HandlerFunc(ret.findMatchingWords))
	mux.Handle("/v1/anagrams", http.HandlerFunc(ret.solveAnagram))
	ret.h = mux
	return ret, nil
}

func (h *Handler) HTTPHandler() http.Handler {
	return h.h
}

func (h *Handler) sendError(w http.ResponseWriter, msg string, code int) {
	ret := map[string]interface{}{
		"error": msg,
		"code":  code,
	}
	b, _ := json.Marshal(ret)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(b)
}

func (h *Handler) synonyms(w http.ResponseWriter, r *http.Request) {
	q, err := synonyms.NewQueryFromParams(r.URL.Query())
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	syns, err := q.Run()
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(syns)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(b)
}

func (h *Handler) findMatchingWords(w http.ResponseWriter, r *http.Request) {
	q, err := findwords.NewQueryFromParams(r.URL.Query())
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := q.Run()
	if err != nil {
		if ok := inputerror.IsInputError(err); ok {
			h.sendError(w, err.Error(), http.StatusBadRequest)
		} else {
			h.sendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(result)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(b)
}

func (h *Handler) solveAnagram(w http.ResponseWriter, r *http.Request) {
	q, err := anagrams.NewQueryFromParams(r.URL.Query())
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := anagrams.Solve(q)
	if err != nil {
		if ok := inputerror.IsInputError(err); ok {
			h.sendError(w, err.Error(), http.StatusBadRequest)
		} else {
			h.sendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(result)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(b)
}
