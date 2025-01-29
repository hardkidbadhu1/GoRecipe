/* 001_create_recipes_collection.js */

db.createCollection("recipes");

// Optionally, create indexes, e.g.:
db.recipes.createIndex({ name: 1 }, { unique: true });