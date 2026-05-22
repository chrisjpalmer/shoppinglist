------- PLAN --------

-- name: GetPlan :one
SELECT * FROM plan LIMIT 1;

-- name: CreatePlan :exec
INSERT INTO plan (plan_data) VALUES (?);

-- name: UpdatePlan :exec
UPDATE plan set plan_data = ? WHERE id = ?;

------- MEALS ------

-- name: GetMeals :many
SELECT 
    id, 
    name, 
    ingredients, 
    recipe_url, 
    preview_image_mode,
    preview_image_url,
    ingredients_image_mode,
    ingredients_image_url
FROM meals
ORDER BY name;

-- name: CreateMeal :one
INSERT INTO meals (
 name,
 ingredients,
 recipe_url,
 preview_image_mode,
 preview_image_url,
 ingredients_image_mode,
 ingredients_image_url
) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id;

-- name: UpdateMeal :exec
UPDATE meals set name = ?,
 ingredients = ?,
 recipe_url = ?,
 preview_image_mode = ?,
 preview_image_url = ?,
 ingredients_image_mode = ?,
 ingredients_image_url = ?
WHERE id = ?;

-- name: DeleteMeal :exec
DELETE FROM meals WHERE id = ?;

-- name: GetMealPreviewImageBytes :one
SELECT preview_image_bytes FROM meals WHERE id = ?;

-- name: GetMealIngredientsImageBytes :one
SELECT ingredients_image_bytes FROM meals WHERE id = ?;

-- name: UpdateMealPreviewImageBytes :exec
UPDATE meals set preview_image_bytes = ? WHERE id = ?;

-- name: UpdateMealIngredientsImageBytes :exec
UPDATE meals set ingredients_image_bytes = ? WHERE id = ?;

------- INGREDIENTS ------

-- name: GetIngredients :many
SELECT id, name, want_override_count, got_count, shopped FROM ingredients
ORDER BY name;

-- name: CreateIngredient :one
INSERT INTO ingredients (name) VALUES (?) RETURNING id;

-- name: UpdateIngredient :exec
UPDATE ingredients set name = ? WHERE id = ?;

-- name: UpdateIngredientWantOverrideCount :exec
UPDATE ingredients set want_override_count = ? WHERE id = ?;

-- name: DeleteIngredient :exec
DELETE FROM ingredients WHERE id = ?;

-- name: ResetIngredientGotCount :exec
UPDATE ingredients SET got_count = 0;

-- name: UpdateIngredientGotCount :exec
UPDATE ingredients SET got_count = ? WHERE id = ?;

-- name: ResetIngredientShopped :exec
UPDATE ingredients SET shopped = FALSE;

-- name: UpdateIngredientShopped :exec
UPDATE ingredients SET shopped = ? WHERE id = ?;

------- INGREDIENT CATEGORIES ------

-- name: GetIngredientCategories :many
SELECT id, name, sort_index FROM ingredient_categories
ORDER BY sort_index ASC;

-- name: CreateIngredientCategory :one
INSERT INTO ingredient_categories (name, sort_index) VALUES (
    ?, 
    (COALESCE((SELECT sort_index FROM ingredient_categories ORDER BY sort_index DESC LIMIT 1), 0) + 1)
) RETURNING id, sort_index;

-- name: UpdateIngredientCategory :exec
UPDATE ingredient_categories SET name = ? WHERE id = ?;

-- name: GetIngredientCategorySortIndex :one
SELECT sort_index FROM ingredient_categories WHERE id = ?;

-- name: UpdateIngredientCategorySortIndex :exec
UPDATE ingredient_categories SET sort_index = ? WHERE id = ?;

-- name: DeleteIngredientCategory :exec
DELETE FROM ingredient_categories WHERE id = ?;

