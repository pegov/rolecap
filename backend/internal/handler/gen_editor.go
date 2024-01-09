package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/pegov/rolecap/backend/internal/extractor"
	"github.com/pegov/rolecap/backend/internal/repo"
	"github.com/pegov/rolecap/backend/pkg/private/editor"
	"github.com/pegov/rolecap/backend/pkg/private/generator"
)

type GenEditorHandler interface {
	Create(c *gin.Context)
	Test(c *gin.Context)
}

type genEditorHandler struct {
	repo          repo.GenRepo
	authExtractor extractor.AuthExtractor
	logger        *slog.Logger
}

func NewGenEditorHandler(
	genRepo repo.GenRepo,
	authExtractor extractor.AuthExtractor,
	logger *slog.Logger,
) GenEditorHandler {
	return &genEditorHandler{genRepo, authExtractor, logger}
}

func (h *genEditorHandler) Test(c *gin.Context) {
	// user, _ := h.authExtractor.GetUser(c)
	// if user == nil {
	// 	c.JSON(401, gin.H{"detail": "Not authenticated"})
	// }

	validator := editor.NewValidator()
	var data editor.EditorRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(422, gin.H{"detail": "Unprocessable entity"})
		return
	}

	report := editor.NewReport()
	meta := validator.ValidateRequestData(&data, report)

	if report.HasError() {
		response := editor.NewEditorErrorResponse(report, nil)
		c.JSON(400, response)
		return
	}

	construct := generator.NewConstruct(data.Body, meta)
	g := generator.NewGenerator()
	result := g.Generate(construct)
	c.JSON(200, result)
}

func (h *genEditorHandler) Create(c *gin.Context) {
	user, _ := h.authExtractor.GetUser(c)
	if user == nil {
		c.JSON(401, gin.H{"detail": "Not authenticated"})
	}

	validator := editor.NewValidator()
	var data editor.EditorRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(422, gin.H{"detail": "Unprocessable entity"})
		return
	}

	report := editor.NewReport()
	meta := validator.ValidateRequestData(&data, report)

	if report.HasError() {
		response := editor.NewEditorErrorResponse(report, nil)
		c.JSON(400, response)
		return
	}

	id, err := h.repo.Create(user, data.Head, data.Body, meta)
	if err != nil {
		c.JSON(500, gin.H{"detail": "Internal server error"})
		return
	}

	response := editor.NewEditorSaveResponse(id)
	c.JSON(200, response)
}
