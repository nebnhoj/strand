package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/nebnhoj/strand/pkg/cache"
	"github.com/nebnhoj/strand/pkg/response"
)

type CacheHandler struct {
	cache cache.Cache
}

func NewCacheHandler(c cache.Cache) *CacheHandler {
	return &CacheHandler{cache: c}
}

// Flush godoc
// @Summary Flush cache
// @Description Clears all entries from the Redis cache (admin only)
// @Tags cache
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /cache [delete]
func (h *CacheHandler) Flush(c fiber.Ctx) error {
	if err := h.cache.Flush(c.Context()); err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}
	return response.Success(c, http.StatusOK, map[string]string{"message": "cache cleared"})
}
