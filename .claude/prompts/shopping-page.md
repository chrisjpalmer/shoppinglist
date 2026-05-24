Implement the shop page of the shopping site (`backend/shopping`) which is meant to show a list of items that the user needs to buy. 

1. new `shopped` column to be added to the database schema managed in `backend/sql/schema.sql`, as well as new queries to be added to `backend/sql/query.sql`
    - add `shopped` column boolean to `ingredients` table
    - new query UpdateIngredientShopped(id, boolean) - sets one ingredient's `shopped` to the value provided
    - new query ResetIngredientShopped - sets all ingredients shopped = false
2. `backend/shopping/page` package gets a new type: `ShopItem`
    - ID
    - Ingredient
    - NeedCount
3. set up shop templates in `backend/shopping/render`
    - accepts a `[]page.ShopItem`
    - table with two columns: 
        - Ingredient: the ingredient to buy 
        - Shopped: a checkbox indicating if the user bought it
    - save and reset button at the bottom ( just like the `got.templ` )
4. set up routes for the `shop` page in `backend/shopping`.
    - create a new file `shop.go`
    - should provide handlers for routes: 
        - `shopping/shop`
        - `shopping/shop/reset
    - should pass a `[]page.ShopItem` to the render
    - the list of ingredients should be used to create to `[]page.ShopItem` but items should be only be included if they meet the following condition: `(ingredientCounts[ing.ID], ing.WantOverrideCount) - ing.GotCount > 0`
    - use `got.go` for inspiration



How to make db changes:
- db queries are managed by `backend/sql/query.sql`
- the db schema is managed by `backend/sql/schema.sql`
- make changes here and then use `dagger generate` ( see codebase skill ) to generate code.

How to make go templ changes.
- Go templ files are located in `backend/shopping/render`. Edit files with the `.templ` extension and leave the `_templ.go` files alone.
- after making changes use `dagger generate` ( see codebase skill ) to generate code.

After generating code, it will be easier to get feedback from the go LSP since the generated code will be present.