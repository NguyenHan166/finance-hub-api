package handlers

import (
	"finance-hub-api/internal/config"
	"finance-hub-api/internal/utils"
	"finance-hub-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	r2Storage *utils.R2Storage
	config    *config.Config
}

func NewUploadHandler(cfg *config.Config) (*UploadHandler, error) {
	// Initialize R2 storage
	r2Storage, err := utils.NewR2Storage(
		cfg.R2.Endpoint,
		cfg.R2.AccessKeyID,
		cfg.R2.SecretAccessKey,
		cfg.R2.BucketName,
		cfg.R2.PublicBaseURL,
	)
	if err != nil {
		return nil, err
	}

	return &UploadHandler{
		r2Storage: r2Storage,
		config:    cfg,
	}, nil
}

// UploadAttachment handles file upload for transaction attachments
func (h *UploadHandler) UploadAttachment(c *gin.Context) {
	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "No file uploaded", err.Error())
		return
	}
	defer file.Close()

	// Validate file
	err = utils.ValidateFile(file, header, h.config.Storage.MaxUploadSize, h.config.Storage.AllowedFileTypes)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid file", err.Error())
		return
	}

	// Upload to R2
	publicURL, err := h.r2Storage.UploadFile(c.Request.Context(), file, header, "attachments")
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload file", err.Error())
		return
	}

	// Return the public URL
	response.SuccessResponse(c, http.StatusOK, "File uploaded successfully", gin.H{
		"url": publicURL,
	})
}

// DeleteAttachment handles file deletion
func (h *UploadHandler) DeleteAttachment(c *gin.Context) {
	// Get URL from request body
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Delete from R2
	err := h.r2Storage.DeleteFile(c.Request.Context(), req.URL)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete file", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "File deleted successfully", nil)
}

// UploadAvatar handles avatar upload
func (h *UploadHandler) UploadAvatar(c *gin.Context) {
	// Get file from form
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "No file uploaded", err.Error())
		return
	}
	defer file.Close()

	// Validate file (only images for avatar)
	allowedTypes := []string{"image/jpeg", "image/png", "image/webp"}
	err = utils.ValidateFile(file, header, h.config.Storage.MaxUploadSize, allowedTypes)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid file", err.Error())
		return
	}

	// Upload to R2
	publicURL, err := h.r2Storage.UploadFile(c.Request.Context(), file, header, "avatars")
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload avatar", err.Error())
		return
	}

	// Return the public URL
	response.SuccessResponse(c, http.StatusOK, "Avatar uploaded successfully", gin.H{
		"avatar_url": publicURL,
	})
}
