package handlers

import (
	"net/http"
	"strconv"

	"gestrym-training/src/nutrition/application/usecases"
	"github.com/gin-gonic/gin"
)

type FoodHandler struct {
	SearchUC  *usecases.SearchFoodsUseCase
	GetByIDUC *usecases.GetFoodByIDUseCase
	ImportUC  *usecases.ImportFoodsWithImagesUseCase
}

func NewFoodHandler(searchUC *usecases.SearchFoodsUseCase, getByIDUC *usecases.GetFoodByIDUseCase, importUC *usecases.ImportFoodsWithImagesUseCase) *FoodHandler {
	return &FoodHandler{
		SearchUC:  searchUC,
		GetByIDUC: getByIDUC,
		ImportUC:  importUC,
	}
}

// SearchFoods godoc
// @Summary      Search foods
// @Description  Retrieve a list of foods filtered by name.
// @Tags         Nutrition
// @Accept       json
// @Produce      json
// @Param        search  query     string  false  "Food name to search"
// @Success      200     {object}  map[string]interface{}
// @Router       /public/foods [get]
func (h *FoodHandler) SearchFoods(c *gin.Context) {
	name := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	foods, total, err := h.SearchUC.Execute(name, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search foods"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  foods,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetFoodByID godoc
// @Summary      Get food by ID
// @Description  Retrieve details of a specific food item.
// @Tags         Nutrition
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Food ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /public/foods/{id} [get]
func (h *FoodHandler) GetFoodByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	food, err := h.GetByIDUC.Execute(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Food not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": food})
}

// ImportFoods godoc
// @Summary      Manual import of foods from USDA
// @Description  Triggers a one-time import of basic food categories from USDA API.
// @Tags         Nutrition
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /public/foods/import [post]
func (h *FoodHandler) ImportFoods(c *gin.Context) {
	err := h.ImportUC.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import foods: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Foods imported successfully"})
}
