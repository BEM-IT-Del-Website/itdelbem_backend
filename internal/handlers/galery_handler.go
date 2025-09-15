package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"bem_be/internal/models"
	"bem_be/internal/services"
	"bem_be/internal/utils"
)

type GaleryHandler struct {
	service *services.GaleryService
}

func NewGaleryHandler(db *gorm.DB) *GaleryHandler {
	return &GaleryHandler{
		service: services.NewGaleryService(db),
	}
}

func (h *GaleryHandler) CreateGalery(c *gin.Context) {
	var galery models.Galery

	if err := c.ShouldBindJSON(&galery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.CreateGalery(&galery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Galery created successfully",
		"data":    galery,
	})
}

func (h *GaleryHandler) UpdateGalery(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var galery models.Galery
	if err := c.ShouldBindJSON(&galery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	galery.ID = uint(id)
	if err := h.service.UpdateGalery(&galery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Galery updated successfully",
		"data":    galery,
	})
}
func (h *GaleryHandler) DeleteGalery(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteGalery(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Photo deleted successfully",
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
		"Berhasil mendapatkan daftar pengumuman",
		metadata,
		galerys,
	)

	c.JSON(http.StatusOK, response)
}
func (h *GaleryHandler) GetGaleryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.service.GetGaleryWithStats(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Galery not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Galery retrieved successfully",
		"data":    result,
	})
}
