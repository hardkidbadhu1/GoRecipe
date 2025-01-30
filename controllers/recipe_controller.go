package controllers

import (
	"GoRecipe/constants"
	"fmt"
	"net/http"

	"GoRecipe/repository"
	"context"
	"github.com/gin-gonic/gin"
)

type SuccessMessage struct {
	Message string `json:"message"`
}

type ApiError struct {
	ErrorMessage string `json:"error_message"`
}

// RecipeController defines the handlers for recipe routes
type RecipeController struct {
	Repo repository.RecipeRepository
}

// NewRecipeController returns a new RecipeController
func NewRecipeController(repo repository.RecipeRepository) *RecipeController {
	return &RecipeController{
		Repo: repo,
	}
}

// CreateRecipe godoc
// @Summary Create a new recipe
// @Description Insert a new recipe record
// @Tags recipes
// @Accept json
// @Produce json
// @Param recipe body repository.Recipe true "Recipe data"
// @Success 201 {object} repository.Recipe
// @Failure 400 {object} ApiError
// @Router /v1/recipes [post]
func (rc *RecipeController) CreateRecipe(c *gin.Context) {
	var recipe repository.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, ApiError{fmt.Sprintf("Invalid request payload - %s", err.Error())})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.MongoTimeout)
	defer cancel()

	created, err := rc.Repo.Create(ctx, &recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiError{err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// GetRecipeByID godoc
// @Summary Get a recipe by ID
// @Description Retrieve a specific recipe using its ID
// @Tags recipes
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 200 {object} repository.Recipe
// @Failure 404 {object} ApiError
// @Router /v1/recipes/{id} [get]
func (rc *RecipeController) GetRecipeByID(c *gin.Context) {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), constants.MongoTimeout)
	defer cancel()

	recipe, err := rc.Repo.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiError{err.Error()})
		return
	}
	if recipe == nil {
		c.JSON(http.StatusNotFound, ApiError{"Recipe not found"})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// GetAllRecipes godoc
// @Summary Get all recipes
// @Description Get a list of all recipes
// @Tags recipes
// @Produce json
// @Success 200 {array} repository.Recipe
// @Router /v1/recipes [get]
func (rc *RecipeController) GetAllRecipes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.MongoTimeout)
	defer cancel()

	recipes, err := rc.Repo.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiError{err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipes)
}

// UpdateRecipe godoc
// @Summary Update a recipe
// @Description Update an existing recipe by ID
// @Tags recipes
// @Accept json
// @Produce json
// @Param id path string true "Recipe ID"
// @Param recipe body repository.Recipe true "Updated recipe data"
// @Success 200 {object} ApiError
// @Failure 400 {object} ApiError
// @Failure 404 {object} ApiError
// @Router /v1/recipes/{id} [put]
func (rc *RecipeController) UpdateRecipe(c *gin.Context) {
	id := c.Param("id")
	var recipe repository.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, ApiError{"Invalid request payload"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.MongoTimeout)
	defer cancel()

	err := rc.Repo.Update(ctx, id, &recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiError{err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessMessage{Message: "Recipe updated successfully"})
}

// DeleteRecipe godoc
// @Summary Delete a recipe
// @Description Remove a recipe by ID
// @Tags recipes
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 204 {object} ApiError
// @Failure 404 {object} ApiError
// @Router /v1/recipes/{id} [delete]
func (rc *RecipeController) DeleteRecipe(c *gin.Context) {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), constants.MongoTimeout)
	defer cancel()

	err := rc.Repo.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, ApiError{err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
