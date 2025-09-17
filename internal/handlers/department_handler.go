package handlers

import (
	"math"
	"net/http"
	"strconv"

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

func (h *DepartmentHandler) GetAllDepartments(c *gin.Context) {
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

    departments, total, err := h.service.GetAllDepartments(perPage, offset, search)
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
        departments,
    )

    c.JSON(http.StatusOK, response)
}

func (h *DepartmentHandler) GetAllDepartmentsGuest(c *gin.Context) {
	// ambil semua data tanpa limit & offset
	departments, err := h.service.GetAllDepartmentsGuest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
		return
	}

	// langsung response tanpa metadata
	response := utils.ResponseHandler(
		"success",
		"Berhasil mendapatkan data",
		departments,
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

	// ambil field manual (biar gak coba bind file ke string)
	department.Name = c.PostForm("name")
	department.ShortName = c.PostForm("short_name")
	department.Vision = c.PostForm("vision")
	department.Mission = c.PostForm("mission")
	department.Value = c.PostForm("value")

	// ambil file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Logo file is required"})
		return
	}

	// kirim ke service
	if err := h.service.CreateDepartment(&department, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Department created successfully",
		"data":    department,
	})
}

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
