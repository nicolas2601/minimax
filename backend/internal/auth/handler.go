package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/nicolas/finanzas/backend/internal/config"
)

type Handler struct {
	svc Service
	cfg *config.Config
}

func NewHandler(svc Service, cfg *config.Config) *Handler {
	return &Handler{svc: svc, cfg: cfg}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}

	user, err := h.svc.Register(req.Email, req.Password, req.DisplayName)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "USER_EXISTS", "message": "El email ya está registrado"}})
			return
		}
		if errors.Is(err, ErrInvalidInput) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "INVALID_INPUT", "message": "Email o contraseña inválidos (mínimo 8 caracteres)"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "Error interno"}})
		return
	}

	result, err := h.svc.Login(user.Email, req.Password, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "Error interno"}})
		return
	}

	setRefreshCookie(c, result.RefreshToken, h.cfg)
	c.JSON(http.StatusCreated, AuthResponse{
		User:         result.User,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}

	result, err := h.svc.Login(req.Email, req.Password, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "INVALID_CREDENTIALS", "message": "Email o contraseña incorrectos"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "Error interno"}})
		return
	}

	setRefreshCookie(c, result.RefreshToken, h.cfg)
	c.JSON(http.StatusOK, AuthResponse{
		User:         result.User,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "NO_REFRESH", "message": "No refresh token"}})
		return
	}

	result, err := h.svc.Refresh(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "INVALID_REFRESH", "message": "Refresh token inválido"}})
		return
	}

	setRefreshCookie(c, result.RefreshToken, h.cfg)
	c.JSON(http.StatusOK, AuthResponse{
		User:         result.User,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")
	if refreshToken != "" {
		_ = h.svc.Logout(refreshToken)
	}

	clearRefreshCookie(c)
	c.Status(http.StatusNoContent)
}

func (h *Handler) Me(c *gin.Context) {
	userAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "No autenticado"}})
		return
	}
	user := userAny.(*User)
	c.JSON(http.StatusOK, UserResponse{User: user})
}

// refreshCookieName is the cookie name used for refresh tokens.
const refreshCookieName = "refresh_token"

// refreshCookiePath scopes the refresh cookie to the auth routes so the
// browser only sends it on /api/v1/auth/* requests, reducing CSRF exposure
// surface and accidental leakage to unrelated paths.
const refreshCookiePath = "/api/v1/auth"

func setRefreshCookie(c *gin.Context, token string, cfg *config.Config) {
	secure := cfg.GinMode != "debug"
	// SetSameSite must be called BEFORE SetCookie because Gin's c.SetSameSite
	// sets the field on c.SetCookie's underlying http.Cookie, and SetCookie
	// copies it onto the response. http.SameSiteLaxMode allows top-level
	// navigation to the auth endpoints (needed for email-link flows) while
	// blocking the cookie on cross-site sub-requests like POST forms and
	// iframes, which mitigates CSRF for state-changing /refresh.
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		refreshCookieName,
		token,
		int(RefreshTokenDuration.Seconds()),
		refreshCookiePath,
		"",
		secure,
		true, // HttpOnly — JS cannot read the refresh token.
	)
}

func clearRefreshCookie(c *gin.Context) {
	// Clear with matching attributes so the browser actually overwrites the
	// cookie. Mismatched Path or SameSite leaves a stale cookie behind.
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshCookieName, "", -1, refreshCookiePath, "", false, true)
}