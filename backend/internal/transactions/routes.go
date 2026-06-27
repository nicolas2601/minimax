package transactions

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts the transaction CRUD + transfer endpoints under
// /transactions. All routes require Bearer auth via requireAuth.
func RegisterRoutes(r *gin.RouterGroup, h *Handler, requireAuth gin.HandlerFunc) {
	g := r.Group("/transactions")
	g.Use(requireAuth)
	{
		g.GET("", h.List)
		g.POST("", h.Create)
		g.GET("/:id", h.Get)
		g.PATCH("/:id", h.Update)
		g.DELETE("/:id", h.Delete)
		g.POST("/transfer", h.Transfer)
	}
}