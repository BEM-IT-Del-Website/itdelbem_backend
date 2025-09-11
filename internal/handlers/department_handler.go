package handlers

import (
	"net/http"
	"strconv"
	"math"
	"fmt"
	"gorm.io/gorm"

	"bem_be/internal/models"
	"bem_be/internal/services"
	"bem_be/internal/utils"
	"github.com/gin-gonic/gin"
)

// DepartmentHandler handles HTTP requests related to departments
type DepartmentHandler struct {
	service *services.DepartmentService
}

// NewDepartmentHandler creates a new department handler
func NewDepartmentHandler(db *gorm.DB) *DepartmentHandler {
	return &DepartmentHandler{
		service: services.NewDepartmentService(db),
	}
}

// GetAllDepartments returns all departments
func (h *DepartmentHandler) GetAllDepartments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

    if page < 1 {
        page = 1
    }
    if perPage < 1 {
        perPage = 10
    }

    offset := (page - 1) * perPage

    // ambil data + total count
    students, total, err := h.service.GetAllDepartments(perPage, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
        return
    }

    totalPages := int(math.Ceil(float64(total) / float64(perPage)))

    // siapkan metadata
    metadata := utils.PaginationMetadata{
        CurrentPage: page,
        PerPage:     perPage,
        TotalItems:  int(total),
        TotalPages:  totalPages,
        Links: utils.PaginationLinks{
            First: fmt.Sprintf("/students?page=1&per_page=%d", perPage),
            Last:  fmt.Sprintf("/students?page=%d&per_page=%d", totalPages, perPage),
        },
    }

    // response dengan metadata
    response := utils.MetadataFormatResponse(
        "success",
        "Berhasil list mendapatkan data",
        metadata,
        students,
    )

    c.JSON(http.StatusOK, response)
}

// GetDepartmentByID returns a department by ID
func (h *DepartmentHandler) GetDepartmentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	stats := c.Query("stats")
	var result interface{}

	if stats == "true" {
		result, err = h.service.GetDepartmentWithStats(uint(id))
	} else {
		result, err = h.service.GetDepartmentByID(uint(id))
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Department retrieved successfully",
		"data":    result,
	})
}

// CreateDepartment creates a new department
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var department models.Department

	if err := c.ShouldBindJSON(&department); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.CreateDepartment(&department); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Department created successfully",
		"data":    department,
	})
}

// UpdateDepartment updates a department
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var department models.Department
	if err := c.ShouldBindJSON(&department); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	department.ID = uint(id)

	if err := h.service.UpdateDepartment(&department); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Department updated successfully",
		"data":    department,
	})
}

// DeleteDepartment deletes a department
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteDepartment(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Department deleted successfully",
	})
} 