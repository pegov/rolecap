package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pegov/rolecap/backend/internal/handler"
)

func SetupGenHandler(
	r *gin.Engine,
	editorHandler handler.GenEditorHandler,
	viewerHandler handler.GenViewerHandler,
) {
	r.POST("/api/gens/editor/create", editorHandler.Create)
	r.POST("/api/gens/editor/test", editorHandler.Test)
	r.GET("/api/gens/:id", viewerHandler.Get)
}
