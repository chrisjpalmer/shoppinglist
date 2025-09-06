CREATE TABLE plan (
  id           INTEGER PRIMARY KEY,
  current_plan boolean,
  plan_data    text NOT NULL -- [[{"category_id": 1, "meal_id": 1}]], [day (0-6)][category/meal assignment]
  next_id      INTEGER,
  prev_id      INTEGER
);

CREATE TABLE meals (
  id           INTEGER PRIMARY KEY,
  name         text    NOT NULL,
  recipe_url   text    NOT NULL,
  ingredients  text    NOT NULL -- [ingredient ids]
);

CREATE TABLE ingredients (
  id INTEGER PRIMARY KEY,
  name text NOT NULL
)