package handlers

import (
	"net/http"
	"strconv"
	"math"
	"gorm.io/gorm"

	"bem_be/internal/models"
	"bem_be/internal/services"
	"bem_be/internal/utils"
	"github.com/gin-gonic/gin"
)

// AssociationHandler handles HTTP requests related to associations
type AssociationHandler struct {
	service *services.AssociationService
}

// NewAssociationHandler creates a new association handler
func NewAssociationHandler(db *gorm.DB) *AssociationHandler {
	return &AssociationHandler{
		service: services.NewAssociationService(db),
	}
}

func (h *AssociationHandler) GetAllAssociations(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
    search := c.Query("name") // pencarian pakai param ?name=

    if page < 1 {
        page = 1
    }
    if perPage < 1 {
        perPage = 10
    }

    offset := (page - 1) * perPage

    associations, total, err := h.service.GetAllAssociations(perPage, offset, search)
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
    }

    response := utils.MetadataFormatResponse(
        "success",
        "Berhasil list mendapatkan data associations",
        metadata,
        associations,
    )

    c.JSON(http.StatusOK, response)
}


func (h *AssociationHandler) GetAllAssociationsGuest(c *gin.Context) {
    // ambil semua data tanpa limit & offset
    associations, err := h.service.GetAllAssociationsGuest()
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
        return
    }

    // langsung response tanpa metadata
    response := utils.ResponseHandler(
        "success",
        "Berhasil mendapatkan data",
        associations,
    )

    c.JSON(http.StatusOK, response)
}

// GetAssociationByID returns a association by ID
func (h *AssociationHandler) GetAssociationByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	stats := c.Query("stats")
	var result interface{}

	if stats == "true" {
		result, err = h.service.GetAssociationWithStats(uint(id))
	} else {
		result, err = h.service.GetAssociationByID(uint(id))
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Association not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Association retrieved successfully",
		"data":    result,
	})
}

// CreateAssociation creates a new association
// CreateAssociation creates a new association
func (h *AssociationHandler) CreateAssociation(c *gin.Context) {
	var association models.Organization

	// ambil field manual (biar gak coba bind file ke string)
	association.Name = c.PostForm("name")
	association.ShortName = c.PostForm("short_name")

	association.CategoryID = 3


	// ambil file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Logo file is required"})
		return
	}

	// kirim ke service
	if err := h.service.CreateAssociation(&association, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Association created successfully",
		"data":    association,
	})
}

// UpdateAssociation updates a association
func (h *AssociationHandler) UpdateAssociation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var association models.Organization
	if err := c.ShouldBindJSON(&association); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	association.ID = uint(id)

	if err := h.service.UpdateAssociation(&association); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Association updated successfully",
		"data":    association,
	})
}

// DeleteAssociation deletes a association
func (h *AssociationHandler) DeleteAssociation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteAssociation(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Association deleted successfully",
	})
} 