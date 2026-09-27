package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) handleListTools(w http.ResponseWriter, r *http.Request) {
	if s.deps.ListTools == nil {
		writeErr(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Tools not available.", nil)
		return
	}
	list, err := s.deps.ListTools(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "TOOLS_FAILED", err.Error(), nil)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (s *Server) handleToolActivity(w http.ResponseWriter, r *http.Request) {
	if s.deps.ToolActivity == nil {
		writeJSON(w, http.StatusOK, []any{})
		return
	}
	writeJSON(w, http.StatusOK, s.deps.ToolActivity())
}

func (s *Server) handleSetToolEnabled(w http.ResponseWriter, r *http.Request) {
	if s.deps.SetToolEnabled == nil {
		writeErr(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Tools not available.", nil)
		return
	}
	var body struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "INVALID_JSON", "invalid body", nil)
		return
	}
	id := mux.Vars(r)["id"]
	if err := s.deps.SetToolEnabled(r.Context(), id, body.Enabled); err != nil {
		writeErr(w, http.StatusBadRequest, "TOOL_UPDATE_FAILED", err.Error(), nil)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": id, "enabled": body.Enabled})
}

func (s *Server) handleTestTool(w http.ResponseWriter, r *http.Request) {
	if s.deps.TestTool == nil {
		writeErr(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Tools not available.", nil)
		return
	}
	var body struct {
		Args map[string]any `json:"args"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "INVALID_JSON", "invalid body", nil)
		return
	}
	result, err := s.deps.TestTool(r.Context(), mux.Vars(r)["id"], body.Args)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "TOOL_TEST_FAILED", err.Error(), nil)
		return
	}
	writeJSON(w, http.StatusOK, result)
}
