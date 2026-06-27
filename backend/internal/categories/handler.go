package categories

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	cat, err := h.svc.Create(userID, req)
	if err != nil {
		if errors.Is(err, ErrInvalidType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_TYPE", "message": "Tipo de categoría inválido"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "Error interno"}})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *Handler) List(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	catType := c.Query("type")
	list, err := h.svc.List(userID, catType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "Error interno"}})
		return
	}
	c.JSON(http.StatusOK, ListResponse{Categories: list})
}

func (h *Handler) Get(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	id, _ := uuid.Parse(c.Param("id"))
	cat, err := h.svc.Get(id, userID)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "Categoría no encontrada"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "Error interno"}})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *Handler) Update(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	id, _ := uuid.Parse(c.Param("id"))
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	cat, err := h.svc.Update(id, userID, req)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "Categoría no encontrada"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "Error interno"}})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *Handler) Delete(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	id, _ := uuid.Parse(c.Param("id"))
	if err := h.svc.Delete(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "Error interno"}})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Seed(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	created, err := h.svc.SeedDefaults(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "Error interno"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"created": created})
}