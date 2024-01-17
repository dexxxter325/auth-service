package handler

import (
	"CRUD_API"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary signUp
/*описание эндпоинта*/
// @Tags auth
/*заголовок для эндпоинта*/
// @Description create account
/*описание внутри эндпоинта*/
// @Accept  json
/*формат request*/
// @Produce  json
/*формат response*/
// @Param input body SignUpRequest true "account info"
// @Router /auth/sign-up [post]
/*эндпоинт нашего хэндлера*/
func (h *Handler) signUp(c *gin.Context) {
	var user CRUD_API.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err in signUp": err.Error()})
		return
	}
	id, err := h.services.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

type SignUpRequest struct { //for swagger
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary signIn
// @Tags auth
// @Description login
// @Accept  json
// @Produce  json
// @Param input body SignInRequest true "credentials"
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var user CRUD_API.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err in signIn": err.Error()})
		return
	}
	accessToken, err := h.services.GenerateAccessToken(user.Username, user.Password, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err in service": err.Error()})
		return
	}
	refreshToken, err := h.services.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err in service": err.Error()})
	}

	h.setRefreshTokenCookie(c, refreshToken) //устанавливаем наш refresh токен в httpOnly куку для предотвращения возможности атак типа XSS.
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

type SignInRequest struct { //for swagger
	ID       int    `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Refresh
// @Tags auth
// @Description refreshtoken pair
// @Accept json
// @Produce json
// @Param refresh_token body RefreshTokenRequest true "Refresh token"
// @Router /auth/refresh-tokens [post]
func (h *Handler) Refresh(c *gin.Context) {
	var request RefreshTokenRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	if request.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}
	newAccessToken, newRefreshToken, err := h.services.GenerateNewTokenPair(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "GenerateNewTokenPair", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// куки-место хранения нашего refresh токена
func (h *Handler) setRefreshTokenCookie(c *gin.Context, refreshToken string) {
	cookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true, // доступна только через HTTP и не доступна из JavaScript для избежания утечки
	}

	// Устанавливаем куку в ответе
	http.SetCookie(c.Writer, cookie) //c.writer-ответ клиенту,а именно-куки
}
