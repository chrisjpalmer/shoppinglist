We need to extend the ingredients table in the `frontend` app to contain the a new column
`ingredient_category_id` that should refer to an id in the `ingredient_categories` table.

1. Extend `ingredients` table to have `ingredient_category_id`:
    - `backend/sql/schema.sql` 
        - Add a new column to ingredients table `ingredient_category_id`
    - `backend/sql/query.sql`
        - CreateIngredient & UpdateIngredient accept `ingredient_category_id`
        - GetIngredients returns `ingredient_category_id` 
2. Extend api server to handle `ingredient_category_id` for ingredient queries 
    - `backend/proto/ingredient.proto` is extended to include `ingredient_category_id`
    - CreateIngredient & UpdateIngredient accept `ingredient_category_id`
    - GetIngredients returns `ingredient_category_id`
3. `frontend/src/routes/ingredients/+page.svelte` extended have extra column "Ingredient Category"
    - `DisplayIngredient` and `DisplayNewIngredient` is extended to contain `ingredient_category_id`
    - new column is inserted before the action column
    - when not in edit mode, shows the ingredient category name
    - when in edit mode, shows a select menu with the various ingredient categories, select is bound to `ingredient_category_id`

How to make db changes:
- db queries are managed by `backend/sql/query.sql`
- the db schema is managed by `backend/sql/schema.sql`
- make changes here and then use `dagger generate` ( see codebase skill ) to generate code.

After generating code, it will be easier to get feedback from the go LSP since the generated code will be present.