CREATE TABLE plan (
  id        INTEGER PRIMARY KEY,
  plan_data text NOT NULL -- [[{"category_id": 1, "meal_id": 1}]], [day (0-6)][category/meal assignment]
);

CREATE TABLE meals (
  id           INTEGER PRIMARY KEY,
  name         text    NOT NULL,
  ingredients  text    NOT NULL -- [ingredient ids]
);

CREATE TABLE ingredients (
  id INTEGER PRIMARY KEY,
  name text NOT NULL
)