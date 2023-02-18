package api

import (
	"encoding/json"
	"net/http"

	"github.com/gotwarlost/crossies/internal/anagrams"
	"github.com/gotwarlost/crossies/internal/findwords"
	"github.com/gotwarlost/crossies/internal/inputerror"
	"github.com/gotwarlost/crossies/internal/synonyms"
	"github.com/naoina/denco"
)

type Handler struct {
	h http.Handler
}

func New() (*Handler, error) {
	mux := denco.NewMux()
	ret := &Handler{}
	h, err := mux.Build([]denco.Handler{
		mux.GET("/v1/synonyms", ret.synonyms),
		mux.GET("/v1/matching-words", ret.findMatchingWords),
		mux.GET("/v1/anagrams", ret.solveAnagram),
	})
	if err != nil {
		return nil, err
	}
	ret.h = h
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

func (h *Handler) synonyms(w http.ResponseWriter, r *http.Request, _ denco.Params) {
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

func (h *Handler) findMatchingWords(w http.ResponseWriter, r *http.Request, _ denco.Params) {
	q, err := findwords.NewQueryFromParams(r.URL.Query())
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := q.Run()
	if err != nil {
		if _, ok := err.(inputerror.InputError); ok {
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

func (h *Handler) solveAnagram(w http.ResponseWriter, r *http.Request, _ denco.Params) {
	q, err := anagrams.NewQueryFromParams(r.URL.Query())
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := anagrams.Solve(q)
	if err != nil {
		if _, ok := err.(inputerror.InputError); ok {
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
