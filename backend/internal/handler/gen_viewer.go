package handler

import (
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pegov/rolecap/backend/internal/extractor"
	"github.com/pegov/rolecap/backend/internal/model"
	"github.com/pegov/rolecap/backend/internal/repo"
)

type GenViewerHandler interface {
	Get(c *gin.Context)
}

type genViewerHandler struct {
	repo          repo.GenRepo
	authExtractor extractor.AuthExtractor
	logger        *slog.Logger
}

func NewGenViewerHandler(
	repo repo.GenRepo,
	authExtractor extractor.AuthExtractor,
	logger *slog.Logger,
) GenViewerHandler {
	return &genViewerHandler{repo, authExtractor, logger}
}

func (h *genViewerHandler) Get(c *gin.Context) {
	user, _ := h.authExtractor.GetUser(c)
	if user == nil {
		c.JSON(401, gin.H{"detail": "Not authenticated"})
	}

	// TODO: check owner
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(404, gin.H{"detail": "Not found"})
		return
	}

	entity, err := h.repo.GetById(id)
	if err != nil {
		c.JSON(500, gin.H{"detail": "Internal server error"})
		return
	}

	if entity == nil {
		c.JSON(404, gin.H{"detail": "Not found"})
		return
	}

	model := model.GenPublicFromGen(entity)
	c.JSON(200, model)
}
