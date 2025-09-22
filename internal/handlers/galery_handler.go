
package handlers

import (
	"bem_be/internal/models"
	"bem_be/internal/services"
	"bem_be/internal/utils"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GaleryHandler struct {
	service *services.GaleryService
}

func NewGaleryHandler(db *gorm.DB) *GaleryHandler {
	return &GaleryHandler{
		service: services.NewGaleryService(db),
	}
}

const (
	maxUploadSize = 5 << 20
	uploadDir     = "uploads/galery"
)

func parseIDParam(c *gin.Context, name string) (uint, bool) {
	idStr := c.Param(name)
	v, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return 0, false
	}
	return uint(v), true
}

func saveImage(c *gin.Context, field string) (string, error) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	file, err := c.FormFile(field)
	if err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file")
	}
	defer src.Close()

	head := make([]byte, 512)
	_, _ = src.Read(head)
	mimeType := http.DetectContentType(head)
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("gagal membaca file")
	}

	allowed := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	if !allowed[mimeType] {
		return "", fmt.Errorf("tipe file tidak didukung (hanya jpg/png/webp)")
	}

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", fmt.Errorf("gagal membuat folder upload")
	}

	name := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	path := filepath.Join(uploadDir, name)

	if err := c.SaveUploadedFile(file, path); err != nil {
		return "", fmt.Errorf("gagal menyimpan file")
	}
	return path, nil
}

func (h *GaleryHandler) CreateGalery(c *gin.Context) {
	ct := c.ContentType()
	var galery models.Galery

	switch {
	case strings.HasPrefix(ct, "multipart/form-data"):
		galery.Title = c.PostForm("title")
		galery.Content = c.PostForm("content")
		path, err := saveImage(c, "image_url")
		if err != nil {
			if err != http.ErrMissingFile {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Gagal memproses file: " + err.Error()})
				return
			}
		} else {
			galery.ImageURL = path
		}

	case strings.HasPrefix(ct, "application/json"):
		if err := c.ShouldBindJSON(&galery); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Body JSON tidak valid"})
			return
		}

	default:
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"status": "error", "message": "Gunakan application/json atau multipart/form-data"})
		return
	}

	if err := h.service.CreateGalery(&galery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Galeri berhasil dibuat",
		"data":    galery,
	})
}
func (h *GaleryHandler) UpdateGalery(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	existing, err := h.service.GetGaleryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Galeri tidak ditemukan"})
		return
	}

	ct := c.ContentType()

	switch {
	case strings.HasPrefix(ct, "multipart/form-data"):
		if v := c.PostForm("title"); v != "" {
			existing.Title = v
		}
		if v := c.PostForm("description"); v != "" {
			existing.Content = v
		}

		// File opsional
		path, err := saveImage(c, "image")
		if err == nil {
			// Hapus file lama jika ada
			if existing.ImageURL != "" {
				_ = os.Remove(existing.ImageURL)
			}
			existing.ImageURL = path
		} else if err != http.ErrMissingFile {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Gagal memproses file: " + err.Error()})
			return
		}

	case strings.HasPrefix(ct, "application/json"):
		var payload models.Galery
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Body JSON tidak valid"})
			return
		}
		if payload.Title != "" {
			existing.Title = payload.Title
		}
		if payload.Content != "" {
			existing.Content = payload.Content
		}
	default:
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"status": "error", "message": "Gunakan application/json atau multipart/form-data"})
		return
	}

	if err := h.service.UpdateGalery(existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Galeri berhasil diperbarui",
		"data":    existing,
	})
}

func (h *GaleryHandler) DeleteGalery(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	gal, _ := h.service.GetGaleryByID(id)

	if err := h.service.DeleteGalery(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	if gal != nil && gal.ImageURL != "" {
		_ = os.Remove(gal.ImageURL)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Foto berhasil dihapus",
	})
}

func (h *GaleryHandler) GetAllGalerys(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	galerys, total, err := h.service.GetAllGalerys(perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	metadata := utils.PaginationMetadata{
		CurrentPage: page,
		PerPage:     perPage,
		TotalItems:  int(total),
		TotalPages:  totalPages,
		Links: utils.PaginationLinks{
			First: fmt.Sprintf("/galerys?page=1&per_page=%d", perPage),
			Last:  fmt.Sprintf("/galerys?page=%d&per_page=%d", totalPages, perPage),
		},
	}

	response := utils.MetadataFormatResponse(
		"success",
		"Berhasil mendapatkan daftar gambar",
		metadata,
		galerys,
	)

	c.JSON(http.StatusOK, response)
}

func (h *GaleryHandler) GetGaleryByID(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	result, err := h.service.GetGaleryWithStats(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Galeri tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Galeri berhasil didapatkan",
		"data":    result,
	})
}