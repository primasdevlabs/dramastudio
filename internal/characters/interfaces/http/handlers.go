package http

import (
	"net/http"
	"strconv"

	"dramastudio/internal/characters/application/services"
	"dramastudio/internal/characters/domain"
	platformhttp "dramastudio/internal/platform/http"
)

type CharactersHandler struct {
	service *services.CharacterService
}

func NewCharactersHandler(service *services.CharacterService) *CharactersHandler {
	return &CharactersHandler{service: service}
}

func (h *CharactersHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/projects/{projectId}/characters", h.list)
	mux.HandleFunc("POST /v1/projects/{projectId}/characters", h.create)
	mux.HandleFunc("GET /v1/projects/{projectId}/characters/{characterId}", h.get)
	mux.HandleFunc("PATCH /v1/projects/{projectId}/characters/{characterId}", h.update)
	mux.HandleFunc("POST /v1/projects/{projectId}/characters/{characterId}/lock", h.lock)
	mux.HandleFunc("DELETE /v1/projects/{projectId}/characters/{characterId}/lock", h.unlock)
	mux.HandleFunc("POST /v1/projects/{projectId}/characters/{characterId}/relationships", h.addRelationship)
	mux.HandleFunc("GET /v1/projects/{projectId}/characters/{characterId}/versions", h.listVersions)
	mux.HandleFunc("GET /v1/projects/{projectId}/characters/{characterId}/versions/{version}", h.getVersion)
	mux.HandleFunc("POST /v1/projects/{projectId}/characters/{characterId}/wardrobe", h.assignWardrobe)
	mux.HandleFunc("GET /v1/projects/{projectId}/characters/{characterId}/wardrobe", h.listWardrobe)
}

// verifyCharacter confirms characterId belongs to the path project —
// prevents hierarchy-confusion access to another project's characters.
func (h *CharactersHandler) verifyCharacter(w http.ResponseWriter, r *http.Request) (*domain.Character, bool) {
	c, err := h.service.GetCharacter(r.Context(), domain.CharacterID(r.PathValue("characterId")))
	if err != nil || c.ProjectID != r.PathValue("projectId") {
		platformhttp.WriteError(w, http.StatusNotFound, "CHARACTER_NOT_FOUND", "Character not found", platformhttp.RequestIDFrom(r), nil)
		return nil, false
	}
	return c, true
}

func (h *CharactersHandler) list(w http.ResponseWriter, r *http.Request) {
	chars, err := h.service.ListCharacters(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	items := make([]CharacterResponse, 0, len(chars))
	for _, c := range chars {
		items = append(items, toCharacterResponse(c))
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": items})
}

func (h *CharactersHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createCharacterReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	c, err := h.service.CreateCharacter(r.Context(), r.PathValue("projectId"), req.Name, req.Role, req.Bio)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, toCharacterResponse(c))
}

func (h *CharactersHandler) get(w http.ResponseWriter, r *http.Request) {
	c, ok := h.verifyCharacter(w, r)
	if !ok {
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, toCharacterResponse(c))
}

func (h *CharactersHandler) update(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.verifyCharacter(w, r); !ok {
		return
	}
	id := domain.CharacterID(r.PathValue("characterId"))
	var req updateCharacterReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	c, err := h.service.UpdateCharacter(r.Context(), id, func(c *domain.Character) error {
		if req.Name != nil {
			c.Name = *req.Name
		}
		if req.Role != nil {
			c.Role = *req.Role
		}
		if req.Bio != nil {
			c.Bio = *req.Bio
		}
		if req.Appearance != nil {
			c.Appearance = *req.Appearance
		}
		if req.Personality != nil {
			c.Personality = *req.Personality
		}
		if req.VoiceProfile != nil {
			c.VoiceProfile = *req.VoiceProfile
		}
		if req.Wardrobe != nil {
			c.Wardrobe = *req.Wardrobe
		}
		return nil
	})
	if err != nil {
		status := http.StatusInternalServerError
		code := "UPDATE_FAILED"
		if err == domain.ErrCharacterNotFound {
			status, code = http.StatusNotFound, "CHARACTER_NOT_FOUND"
		} else if err == domain.ErrCharacterLocked {
			status, code = http.StatusConflict, "CHARACTER_LOCKED"
		}
		platformhttp.WriteError(w, status, code, err.Error(), platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, toCharacterResponse(c))
}

func (h *CharactersHandler) lock(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.verifyCharacter(w, r); !ok {
		return
	}
	if err := h.service.Lock(r.Context(), domain.CharacterID(r.PathValue("characterId"))); err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "CHARACTER_NOT_FOUND", "Character not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CharactersHandler) unlock(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.verifyCharacter(w, r); !ok {
		return
	}
	if err := h.service.Unlock(r.Context(), domain.CharacterID(r.PathValue("characterId"))); err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "CHARACTER_NOT_FOUND", "Character not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CharactersHandler) addRelationship(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.verifyCharacter(w, r); !ok {
		return
	}
	var req relationshipReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	c, err := h.service.AddRelationship(r.Context(), domain.CharacterID(r.PathValue("characterId")), domain.Relationship{
		TargetCharacterID: domain.CharacterID(req.TargetCharacterID),
		RelationshipType:  req.RelationshipType,
		Description:       req.Description,
	})
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, toCharacterResponse(c))
}

func (h *CharactersHandler) listVersions(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.verifyCharacter(w, r); !ok {
		return
	}
	vs, err := h.service.ListVersions(r.Context(), domain.CharacterID(r.PathValue("characterId")))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": vs})
}

func (h *CharactersHandler) getVersion(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.verifyCharacter(w, r); !ok {
		return
	}
	v, err := strconv.Atoi(r.PathValue("version"))
	if err != nil || v < 1 {
		platformhttp.WriteError(w, http.StatusBadRequest, "INVALID_VERSION", "Version must be a positive integer", platformhttp.RequestIDFrom(r), nil)
		return
	}
	cv, err := h.service.GetVersion(r.Context(), domain.CharacterID(r.PathValue("characterId")), v)
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "VERSION_NOT_FOUND", "Character version not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, cv)
}

func (h *CharactersHandler) assignWardrobe(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.verifyCharacter(w, r); !ok {
		return
	}
	var req wardrobeReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	wa, err := h.service.AssignWardrobe(r.Context(), domain.CharacterID(r.PathValue("characterId")),
		req.EpisodeID, req.SceneID, req.Items, req.ChangeEvent)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, wa)
}

func (h *CharactersHandler) listWardrobe(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.verifyCharacter(w, r); !ok {
		return
	}
	episodeID := r.URL.Query().Get("episode_id")
	items, err := h.service.ListWardrobe(r.Context(), domain.CharacterID(r.PathValue("characterId")), episodeID)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": items})
}
