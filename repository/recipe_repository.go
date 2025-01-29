package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Recipe represents the recipe model stored in MongoDB
type Recipe struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `bson:"name" json:"name"`
	Ingredients  []string           `bson:"ingredients" json:"ingredients"`
	Instructions string             `bson:"instructions" json:"instructions"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type RecipeRepository interface {
	Create(ctx context.Context, recipe *Recipe) (*Recipe, error)
	GetByID(ctx context.Context, id string) (*Recipe, error)
	GetAll(ctx context.Context) ([]*Recipe, error)
	Update(ctx context.Context, id string, recipe *Recipe) error
	Delete(ctx context.Context, id string) error
}

// recipeRepository is the implementation of RecipeRepository
type recipeRepository struct {
	collection *mongo.Collection
}

// NewRecipeRepository returns a new instance of RecipeRepository
func NewRecipeRepository(db *mongo.Database, collectionName string) RecipeRepository {
	return &recipeRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *recipeRepository) Create(ctx context.Context, recipe *Recipe) (*Recipe, error) {
	recipe.ID = primitive.NewObjectID()
	recipe.CreatedAt = time.Now()
	recipe.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, recipe)
	if err != nil {
		return nil, fmt.Errorf("failed to insert recipe: %w", err)
	}
	return recipe, nil
}

func (r *recipeRepository) GetByID(ctx context.Context, id string) (*Recipe, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	var recipe Recipe
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&recipe)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe by id: %w", err)
	}
	return &recipe, nil
}

func (r *recipeRepository) GetAll(ctx context.Context) ([]*Recipe, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find recipes: %w", err)
	}
	defer cursor.Close(ctx)

	var recipes []*Recipe
	for cursor.Next(ctx) {
		var recipe Recipe
		if err := cursor.Decode(&recipe); err != nil {
			return nil, fmt.Errorf("failed to decode recipe: %w", err)
		}
		recipes = append(recipes, &recipe)
	}

	return recipes, nil
}

func (r *recipeRepository) Update(ctx context.Context, id string, recipe *Recipe) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	recipe.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":         recipe.Name,
			"ingredients":  recipe.Ingredients,
			"instructions": recipe.Instructions,
			"updated_at":   recipe.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update, options.Update().SetUpsert(false))
	if err != nil {
		return fmt.Errorf("failed to update recipe: %w", err)
	}
	return nil
}

func (r *recipeRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete recipe: %w", err)
	}
	return nil
}
