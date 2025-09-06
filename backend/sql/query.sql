------- PLAN --------

-- name: GetCurrentPlan :many
SELECT id FROM plan WHERE current_plan = true;

-- name: GetPlan :one
SELECT * FROM plan WHERE id = ?;

-- name: CreatePlan :exec
INSERT INTO plan (plan_data) VALUES (?);

-- name: UpdatePlan :exec
UPDATE plan set plan_data = ? WHERE id = ?;

-- name: UnmarkPlansAsCurrent :exec
UPDATE plan SET current_plan = false;

-- name: MarkPlanAsCurrent :exec
UPDATE plan SET curent_plan = true WHERE id = ?;

------- MEALS ------

-- name: GetMeals :many
SELECT id, name, ingredients, recipe_url FROM meals
ORDER BY name;

-- name: CreateMeal :one
INSERT INTO meals (name, ingredients, recipe_url) VALUES (?, ?, ?) RETURNING id;

-- name: UpdateMeal :exec
UPDATE meals set name = ?, ingredients = ?, recipe_url = ? WHERE id = ?;

-- name: DeleteMeal :exec
DELETE FROM meals WHERE id = ?;

------- INGREDIENTS ------

-- name: GetIngredients :many
SELECT id, name FROM ingredients
ORDER BY name;

-- name: CreateIngredient :one
INSERT INTO ingredients (name) VALUES (?) RETURNING id;

-- name: UpdateIngredient :exec
UPDATE ingredients set name = ? WHERE id = ?;

-- name: DeleteIngredient :exec
DELETE FROM ingredients WHERE id = ?;