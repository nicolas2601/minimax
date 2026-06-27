package transactions

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	var dto CreateRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	tx, err := h.svc.Create(userID, dto.ToServiceCreate())
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, tx)
}

func (h *Handler) List(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	var q ListFilterQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	list, err := h.svc.List(userID, q.ToListFilter())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "Error interno"}})
		return
	}
	c.JSON(http.StatusOK, ListResponse{Transactions: list})
}

func (h *Handler) Get(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	id, _ := uuid.Parse(c.Param("id"))
	tx, err := h.svc.Get(id, userID)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, tx)
}

func (h *Handler) Update(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	id, _ := uuid.Parse(c.Param("id"))
	var dto UpdateRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	tx, err := h.svc.Update(id, userID, dto.ToServiceUpdate())
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, tx)
}

func (h *Handler) Delete(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	id, _ := uuid.Parse(c.Param("id"))
	if err := h.svc.Delete(id, userID); err != nil {
		writeServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Transfer(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	var dto TransferRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	res, err := h.svc.Transfer(userID, dto.ToServiceTransfer())
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, TransferResponse{Transfer: *res})
}

func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrTransactionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "Transacción no encontrada"}})
	case errors.Is(err, ErrInvalidType):
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_TYPE", "message": "Tipo de transacción inválido"}})
	case errors.Is(err, ErrAccountNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "ACCOUNT_NOT_FOUND", "message": "Cuenta no encontrada o no pertenece al usuario"}})
	case errors.Is(err, ErrCategoryNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "CATEGORY_NOT_FOUND", "message": "Categoría no encontrada o no pertenece al usuario"}})
	case errors.Is(err, ErrInvalidAmount):
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_AMOUNT", "message": "El monto debe ser mayor a cero"}})
	case errors.Is(err, ErrSameAccount):
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "SAME_ACCOUNT", "message": "La cuenta de origen y destino deben ser distintas"}})
	case errors.Is(err, ErrCurrencyMismatch):
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "CURRENCY_MISMATCH", "message": "Las cuentas deben compartir la misma moneda"}})
	default:
		// Date parse errors fall through here. We sniff the message so the
		// handler stays decoupled from the service's error package.
		msg := err.Error()
		if hasPrefix(msg, "invalid date format") {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_DATE", "message": msg}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "Error interno"}})
	}
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}