package handlers

import (
	"bem_be/internal/models"
	"bem_be/internal/services"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// VisiMisiHandler menangani request HTTP terkait berita
type VisiMisiHandler struct {
	service *services.VisiMisiService
}

// NewVisiMisiHandler membuat handler berita baru
func NewVisiMisiHandler(db *gorm.DB) *VisiMisiHandler {
	return &VisiMisiHandler{
		service: services.NewVisiMisiService(db),
	}
}

// GetAllVisiMisi mengembalikan semua berita dengan pagination
func (h *VisiMisiHandler) GetAllVisiMisi(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	visimisiList, total, err := h.service.GetAllVisiMisi(perPage, offset)
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
		"message": "Berhasil mendapatkan daftar berita",
		"metadata": gin.H{
			"current_page": page,
			"per_page":     perPage,
			"total_items":  total,
			"total_pages":  totalPages,
		},
		"data": visimisiList,
	})
}

// GetVisiMisiByID mengembalikan berita berdasarkan ID
func (h *VisiMisiHandler) GetVisiMisiByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	visimisi, err := h.service.GetVisiMisiByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berita berhasil didapatkan",
		"data":    visimisi,
	})
}

// CreateVisiMisi membuat berita baru (dengan unggahan file opsional)
func (h *VisiMisiHandler) CreateVisiMisi(c *gin.Context) {
	var visimisi models.Period

	visimisi.Vision = c.PostForm("visi")
	visimisi.Mission = c.PostForm("misi")

	if err := h.service.CreateVisiMisi(&visimisi); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Berita berhasil dibuat",
		"data":    visimisi,
	})
}

// UpdateVisiMisi memperbarui berita yang ada.
func (h *VisiMisiHandler) UpdateVisiMisi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	existingVisiMisi, err := h.service.GetVisiMisiByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
		return
	}

	existingVisiMisi.Vision = c.PostForm("title")
	existingVisiMisi.Mission = c.PostForm("content")

	if err := h.service.UpdateVisiMisi(existingVisiMisi); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berita berhasil diperbarui",
		"data":    existingVisiMisi,
	})
}

// DeleteVisiMisi menghapus sebuah berita
func (h *VisiMisiHandler) DeleteVisiMisi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	if err := h.service.DeleteVisiMisi(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berita berhasil dihapus",
	})
}

// RestoreVisiMisi menangani permintaan untuk memulihkan berita yang telah di-soft-delete.
func (h *VisiMisiHandler) RestoreVisiMisi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	restoredVisiMisi, err := h.service.RestoreVisiMisi(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berita berhasil dipulihkan",
		"data":    restoredVisiMisi,
	})
}
