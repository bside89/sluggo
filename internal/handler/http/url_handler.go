package httphandler

import (
	"errors"
	"net/http"
	"net/url"

	"sluggo/internal/domain"
	"sluggo/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// URLHandler exposes the HTTP interface for URL shortening operations.
type URLHandler struct {
	uc *usecase.URLUseCase
}

// NewURLHandler creates a URLHandler with the given use case.
func NewURLHandler(uc *usecase.URLUseCase) *URLHandler {
	return &URLHandler{uc: uc}
}

func init() {
	// Register custom validation functions for request binding.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("http_url", validateHTTPURL)
	}
}

type shortenRequest struct {
	URL string `json:"url" binding:"required,http_url" example:"https://www.example.com/some/long/path"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url" example:"http://localhost:8080/abc123"`
}

// Shorten godoc
// @Summary 	Shorten a long URL
// @Description Accepts a long URL and returns a shortened version.
// @Tags 		URL
// @Accept 		json
// @Produce 	json
// @Param 		request body shortenRequest true "URL to shorten"
// @Success 	201 {object} shortenResponse
// @Failure 	400 {object} map[string]string "Bad Request"
// @Failure 	500 {object} map[string]string "Internal Server Error"
// @Router 		/shorten [post]
func (h *URLHandler) Shorten(c *gin.Context) {
	var req shortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": shortenRequestBindError(err)})
		return
	}

	shortURL, err := h.uc.ShortenURL(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shortenResponse{ShortURL: shortURL})
}

// Redirect godoc
// @Summary 	Redirect to the original URL
// @Description Given a short hash, redirects to the associated long URL.
// @Tags 		URL
// @Param 		hash path string true "Short URL hash"
// @Success 	302 "Redirect"
// @Failure 	404 {object} map[string]string "Not Found"
// @Failure 	500 {object} map[string]string "Internal Server Error"
// @Router 		/{hash} [get]
func (h *URLHandler) Redirect(c *gin.Context) {
	hash := c.Param("hash")

	longURL, err := h.uc.ResolveURL(hash)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "short url not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Redirect(http.StatusFound, longURL)
}

// shortenRequestBindError translates JSON binding and validation errors into user-friendly messages.
func shortenRequestBindError(err error) string {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return "invalid request body"
	}

	for _, validationErr := range validationErrs {
		if validationErr.StructField() != "URL" {
			continue
		}

		switch validationErr.Tag() {
		case "required":
			return "field 'url' is required"
		case "http_url":
			return "field 'url' must be a valid URL"
		}
	}

	return "invalid request body"
}

// validateHTTPURL is a custom validator function that checks if a string is a valid HTTP or HTTPS URL.
func validateHTTPURL(fl validator.FieldLevel) bool {
	rawURL := fl.Field().String()
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	return parsedURL.Host != ""
}
