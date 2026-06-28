package recurring

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

// RegisterRoutes binds CRUD + manual run + list runs.
func RegisterRoutes(rg *gin.RouterGroup, h *Handler, requireUserID gin.HandlerFunc) {
	g := rg.Group("/recurring-rules", requireUserID)
	g.GET("", h.list)
	g.POST("", h.create)
	g.GET(":id", h.get)
	g.PATCH(":id", h.update)
	g.DELETE(":id", h.delete)
	g.POST(":id/run-now", h.runNow)
	g.GET(":id/runs", h.listRuns)

	// Cron-style endpoint for the recurring worker. Authenticated so we
	// can audit who triggered generation manually.
	rg.POST("/recurring/generate-today", requireUserID, h.generateToday)
}

const userIDContextKey = "userID"

func (h *Handler) userID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get(userIDContextKey)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "MISSING_AUTH_CONTEXT", "message": "userID missing from request context"},
		})
		return uuid.Nil, false
	}
	s, ok := v.(string)
	if !ok || s == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INVALID_AUTH_CONTEXT", "message": "userID has wrong type"},
		})
		return uuid.Nil, false
	}
	id, err := uuid.Parse(s)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{"code": "BAD_USER_ID", "message": "userID is not a valid UUID"},
		})
		return uuid.Nil, false
	}
	return id, true
}

func (h *Handler) list(c *gin.Context) {
	uid, ok := h.userID(c)
	if !ok {
		return
	}
	rules, err := h.svc.List(uid)
	if err != nil {
		serverError(c, err)
		return
	}
	out := make([]*RuleDTO, 0, len(rules))
	for _, r := range rules {
		out = append(out, r.ToDTO())
	}
	c.JSON(http.StatusOK, gin.H{"recurring_rules": out})
}

func (h *Handler) create(c *gin.Context) {
	uid, ok := h.userID(c)
	if !ok {
		return
	}
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "JSON inválido: "+err.Error())
		return
	}
	rule, err := h.svc.Create(uid, req)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, rule.ToDTO())
}

func (h *Handler) get(c *gin.Context) {
	uid, ok := h.userID(c)
	if !ok {
		return
	}
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	rule, err := h.svc.Get(id, uid)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, rule.ToDTO())
}

func (h *Handler) update(c *gin.Context) {
	uid, ok := h.userID(c)
	if !ok {
		return
	}
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "JSON inválido: "+err.Error())
		return
	}
	rule, err := h.svc.Update(id, uid, req)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, rule.ToDTO())
}

func (h *Handler) delete(c *gin.Context) {
	uid, ok := h.userID(c)
	if !ok {
		return
	}
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Delete(id, uid); err != nil {
		mapServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) runNow(c *gin.Context) {
	uid, ok := h.userID(c)
	if !ok {
		return
	}
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	txID, err := h.svc.RunNow(id, uid)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"transaction_id": txID})
}

func (h *Handler) listRuns(c *gin.Context) {
	uid, ok := h.userID(c)
	if !ok {
		return
	}
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	runs, err := h.svc.ListRuns(id, uid, 50)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	out := make([]*RunDTO, 0, len(runs))
	for _, r := range runs {
		out = append(out, r.ToDTO())
	}
	c.JSON(http.StatusOK, gin.H{"runs": out})
}

func (h *Handler) generateToday(c *gin.Context) {
	stats, err := h.svc.GenerateToday()
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, stats.AsDTO())
}

func parseID(c *gin.Context, key string) (uuid.UUID, bool) {
	raw := c.Param(key)
	id, err := uuid.Parse(raw)
	if err != nil {
		badRequest(c, "INVALID_ID", "id must be a UUID")
		return uuid.Nil, false
	}
	return id, true
}

func badRequest(c *gin.Context, code, msg string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"error": gin.H{"code": code, "message": msg},
	})
}

func serverError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
	})
}

func mapServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrRuleNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "recurring rule not found"},
		})
	case errors.Is(err, ErrInvalidAmount),
		errors.Is(err, ErrInvalidFrequency),
		errors.Is(err, ErrInvalidType),
		errors.Is(err, ErrInvalidInterval),
		errors.Is(err, ErrInvalidStartDate),
		errors.Is(err, ErrEndBeforeStart):
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
	default:
		serverError(c, err)
	}
}
