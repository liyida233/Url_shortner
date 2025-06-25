package handler

import (
	"URL_shortener/shortener"
	"URL_shortener/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	store store.Store
}

func NewHandler(s store.Store) *Handler {
	return &Handler{store: s}
}

type UrlCreationRequest struct {
	LongUrl     string `json:"long_url" binding:"required"`
	UserId      string `json:"user_id" binding:"required"`
	CustomAlias string `json:"custom_alias"`
}

func (h *Handler) CreateShortUrl(c *gin.Context) {
	var req UrlCreationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl := req.CustomAlias
	if shortUrl == "" {
		shortUrl = shortener.GenerateShortLink(req.LongUrl, req.UserId)
	} else {
		if h.store.Exists(shortUrl) {
			c.JSON(http.StatusConflict, gin.H{"error": "Custom alias is already in use."})
			return
		}
	}

	if err := h.store.Save(shortUrl, req.LongUrl, req.UserId, req.CustomAlias); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save short URL"})
		return
	}

	host := "http://localhost:9808/"
	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": host + shortUrl,
	})
}

func (h *Handler) HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	longUrl, err := h.store.Get(shortUrl)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// 类型断言：如果支持访问统计接口就调用
	if tracker, ok := h.store.(store.VisitTracker); ok {
		tracker.IncrementVisitCount(shortUrl)
	}

	c.Redirect(http.StatusFound, longUrl)
}
