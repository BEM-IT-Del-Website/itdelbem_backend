package handlers

import (
	"net/http"
	"strconv"
	"gorm.io/gorm"

	"bem_be/internal/services"
	"github.com/gin-gonic/gin"
	
)

// OrganizationHandler handles HTTP requests related to associations
type OrganizationHandler struct {
	service *services.OrganizationService
}

func NewOrganizationHandler(db *gorm.DB) *OrganizationHandler {
	return &OrganizationHandler{
		service: services.NewOrganizationService(db),
	}
}

// GetOrganizationByID returns a association by ID
func (h *OrganizationHandler) GetOrganizationByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	stats := c.Query("stats")
	var result interface{}

	if stats == "true" {
		result, err = h.service.GetOrganizationWithStats(uint(id))
	} else {
		result, err = h.service.GetOrganizationByID(uint(id))
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Organization retrieved successfully",
		"data":    result,
	})
}