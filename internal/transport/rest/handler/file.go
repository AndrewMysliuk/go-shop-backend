package handler

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/AndrewMislyuk/go-shop-backend/internal/transport/rest/domain"
	"github.com/gin-gonic/gin"
)

const (
	maxUploadSize = 5 << 20 // 5 megabytes
)

var (
	imageTypes = map[string]interface{}{
		"image/jpeg": nil,
		"image/png":  nil,
		"image/jpg":  nil,
		"image/webp": nil,
	}
)

type UploadedImageURL struct {
	URL string `json:"image_url"`
}

// @Summary Upload image
// @Security ApiKeyAuth
// @Tags Upload image
// @Description upload image
// @ID file-upload-image
// @Accept json
// @Produce json
// @Param productId formData string true "productId"
// @Param file formData file true "file"
// @Success 200 {object} UploadedImageURL
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/file/upload [post]
func (h *Handler) uploadImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	err := c.Request.ParseForm()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	productId := c.PostForm("productId")
	if productId == "" {
		newErrorResponse(c, http.StatusBadRequest, "select product id")

		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	defer file.Close()

	buffer := make([]byte, header.Size)

	if _, err := file.Read(buffer); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	contentType := http.DetectContentType(buffer)

	// Validate File Type
	if _, ex := imageTypes[contentType]; !ex {
		newErrorResponse(c, http.StatusBadRequest, "file type is not supported")

		return
	}

	filename := header.Filename

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to create temp file")

		return
	}

	if _, err := io.Copy(f, bytes.NewReader(buffer)); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to write chunk to temp file")

		return
	}

	f.Close()

	url, err := h.fileService.Upload(domain.File{
		ProductId:   productId,
		Type:        domain.Image,
		ContentType: contentType,
		Name:        filename,
		Size:        header.Size,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, UploadedImageURL{
		URL: url,
	})
}
