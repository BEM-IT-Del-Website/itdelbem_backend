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

// ClubHandler handles HTTP requests related to clubs
type ClubHandler struct {
	service *services.ClubService
}

// NewClubHandler creates a new club handler
func NewClubHandler(db *gorm.DB) *ClubHandler {
	return &ClubHandler{
		service: services.NewClubService(db),
	}
}

// GetAllClubs returns all clubs
func (h *ClubHandler) GetAllClubs(c *gin.Context) {
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

    clubs, total, err := h.service.GetAllClubs(perPage, offset, search)
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
        "Berhasil list mendapatkan data clubs",
        metadata,
        clubs,
    )

    c.JSON(http.StatusOK, response)
}


func (h *ClubHandler) GetAllClubsGuest(c *gin.Context) {
    // ambil semua data tanpa limit & offset
    clubs, err := h.service.GetAllClubsGuest()
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
        return
    }

    // langsung response tanpa metadata
    response := utils.ResponseHandler(
        "success",
        "Berhasil mendapatkan data",
        clubs,
    )

    c.JSON(http.StatusOK, response)
}

// GetClubByID returns a club by ID
func (h *ClubHandler) GetClubByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	stats := c.Query("stats")
	var result interface{}

	if stats == "true" {
		result, err = h.service.GetClubWithStats(uint(id))
	} else {
		result, err = h.service.GetClubByID(uint(id))
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Club retrieved successfully",
		"data":    result,
	})
}

// CreateClub creates a new club
func (h *ClubHandler) CreateClub(c *gin.Context) {
	var club models.Organization

	// ambil field manual (biar gak coba bind file ke string)
	club.Name = c.PostForm("name")
	club.ShortName = c.PostForm("short_name")
	
	club.CategoryID = 1


	// ambil file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Logo file is required"})
		return
	}

	// kirim ke service
	if err := h.service.CreateClub(&club, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Club created successfully",
		"data":    club,
	})
}

// UpdateClub updates a club
func (h *ClubHandler) UpdateClub(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var club models.Organization
	if err := c.ShouldBindJSON(&club); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	club.ID = uint(id)

	if err := h.service.UpdateClub(&club); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Club updated successfully",
		"data":    club,
	})
}

// DeleteClub deletes a club
func (h *ClubHandler) DeleteClub(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteClub(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Club deleted successfully",
	})
} 