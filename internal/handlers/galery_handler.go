package handlers

import (
	"bem_be/internal/models"
	"bem_be/internal/services"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
func parseIDParam(value string) *uint {
	if value == "" {
		return nil
	}
	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return nil
	}
	u := uint(val)
	return &u
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
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var totalPages int
	if perPage > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(perPage)))
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mendapatkan daftar galeri",
		"metadata": gin.H{
			"current_page": page,
			"per_page":     perPage,
			"total_items":  total,
			"total_pages":  totalPages,
		},
		"data": galerys,
	})
}
func (h *GaleryHandler) GetGaleryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	galery, err := h.service.GetGaleryByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Galeri berhasil didapatkan",
		"data":    galery,
	})
}
func (h *GaleryHandler) CreateGalery(c *gin.Context) {
	var galery models.Galery
	galery.Title = c.PostForm("title")
	galery.Content = c.PostForm("content")
	galery.Category = c.PostForm("category")
	userID := c.MustGet("userID").(uint)
	galery.UserID = userID
	file, err := c.FormFile("image")
	if err == nil {
		uploadPath := "uploads/galery"
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Tidak dapat membuat folder unggahan"})
			return
		}
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
		filePath := filepath.Join(uploadPath, fileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan file"})
			return
		}
		galery.ImageURL = filePath
	} else if err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Gagal memproses file: " + err.Error()})
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
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}
	existingGalery, err := h.service.GetGaleryByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
		return
	}
	existingGalery.Title = c.PostForm("title")
	existingGalery.Content = c.PostForm("content")
	existingGalery.Category = c.PostForm("category")
	file, err := c.FormFile("image")
	if err == nil {
		if existingGalery.ImageURL != "" {
			_ = os.Remove(existingGalery.ImageURL)
		}
		uploadPath := "uploads/galery"
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Tidak dapat membuat folder unggahan"})
			return
		}
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
		filePath := filepath.Join(uploadPath, fileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan file"})
			return
		}
		existingGalery.ImageURL = filePath
	} else if err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Gagal memproses file: " + err.Error()})
		return
	}
	if err := h.service.UpdateGalery(existingGalery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Galeri berhasil diperbarui",
		"data":    existingGalery,
	})
}
func (h *GaleryHandler) DeleteGalery(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}
	gal, _ := h.service.GetGaleryByID(uint(id))
	if err := h.service.DeleteGalery(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	if gal != nil && gal.ImageURL != "" {
		_ = os.Remove(gal.ImageURL)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Galeri berhasil dihapus",
	})
}
