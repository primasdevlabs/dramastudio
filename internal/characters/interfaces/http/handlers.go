package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"dramastudio/internal/characters/application/services"
	"dramastudio/internal/characters/domain"
	platformHTTP "dramastudio/internal/platform/http"
)

type CharactersHandler struct {
	service *services.CharacterService
}

func NewCharactersHandler(service *services.CharacterService) *CharactersHandler {
	return &CharactersHandler{service: service}
}

type createCharacterReq struct {
	ProjectID string `json:"project_id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Bio       string `json:"bio"`
}

func (h *CharactersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/characters")
	path = strings.TrimPrefix(path, "/")

	if path == "" || path == "/" {
		switch r.Method {
		case http.MethodGet:
			h.listCharacters(w, r)
		case http.MethodPost:
			h.createCharacter(w, r)
		default:
			platformHTTP.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "", nil)
		}
		return
	}

	parts := strings.Split(path, "/")
	charID := domain.CharacterID(parts[0])

	if len(parts) == 1 {
		if r.Method == http.MethodGet {
			h.getCharacter(w, r, charID)
			return
		}
	} else if len(parts) == 2 && parts[1] == "wardrobe" {
		if r.Method == http.MethodPost {
			h.updateWardrobe(w, r, charID)
			return
		}
	}

	platformHTTP.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found", "", nil)
}

func (h *CharactersHandler) listCharacters(w http.ResponseWriter, r *http.Request) {
	chars, err := h.service.ListCharacters(r.Context())
	if err != nil {
		platformHTTP.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), "", nil)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"characters": chars})
}

func (h *CharactersHandler) createCharacter(w http.ResponseWriter, r *http.Request) {
	var req createCharacterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
		return
	}
	c, err := h.service.CreateCharacter(r.Context(), req.ProjectID, req.Name, req.Role, req.Bio)
	if err != nil {
		platformHTTP.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error(), "", nil)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusCreated, c)
}

func (h *CharactersHandler) getCharacter(w http.ResponseWriter, r *http.Request, id domain.CharacterID) {
	c, err := h.service.GetCharacter(r.Context(), id)
	if err != nil {
		platformHTTP.WriteError(w, http.StatusNotFound, "CHARACTER_NOT_FOUND", err.Error(), "", nil)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusOK, c)
}

func (h *CharactersHandler) updateWardrobe(w http.ResponseWriter, r *http.Request, id domain.CharacterID) {
	var wardrobe domain.Wardrobe
	if err := json.NewDecoder(r.Body).Decode(&wardrobe); err != nil {
		platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
		return
	}
	c, err := h.service.UpdateWardrobe(r.Context(), id, wardrobe)
	if err != nil {
		platformHTTP.WriteError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error(), "", nil)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusOK, c)
}
