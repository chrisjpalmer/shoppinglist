We need to create a new entity called "ingredient categories" changes are required to both 
the frontend/backend. The following changes are required

1. New Table ( ingredient_categories )
    1. Table (added to `backend/sql/schema.sql`):
        - id primary key
        - name text NOT NULL
        - sort_index int NOT NULL
    2. New queries added to `backend/sql/queries.sql`:
        -- name: GetIngredients :many
        -- name: CreateIngredient :one
        -- name: UpdateIngredient :exec
        -- name: DeleteIngredient :exec
2. new protos
    - rpc GetIngredients(GetIngredientCategoriesRequest) returns (GetIngredientCategoriesResponse);
    - rpc CreateIngredientCategory(CreateIngredientCategoryRequest) returns (CreateIngredientCategoryResponse);
    - rpc UpdateIngredientCategory(UpdateIngredientCategoryRequest) returns (UpdateIngredientCategoryResponse);
    - rpc DeleteIngredientCategory(DeleteIngredientCategoryRequest) returns (DeleteIngredientCategoryResponse);
3. new `backend/api/ingredient_categories.go` file which is similar to `backend/api/ingredients.go` (essentially maps reqs/resp between proto and database)
3. new svelte page at `frontend/src/routes/ingredient-categories/+page.svelte`
    - heavily based on implementation of ingredients page
    - should have two columns:
        - Name
        - Action
4. new menu item added to `frontend/src/routes/Nav.svelte` to permit navigation to this page

How to make db changes:
- db queries are managed by `backend/sql/query.sql`
- the db schema is managed by `backend/sql/schema.sql`
- make changes here and then use `dagger generate` ( see codebase skill ) to generate code.

How to make go templ changes.
- Go templ files are located in `backend/shopping/render`. Edit files with the `.templ` extension and leave the `_templ.go` files alone.
- after making changes use `dagger generate` ( see codebase skill ) to generate code.

After generating code, it will be easier to get feedback from the go LSP since the generated code will be present.