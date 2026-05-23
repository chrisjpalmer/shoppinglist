CREATE TABLE plan (
  id        INTEGER PRIMARY KEY,
  plan_data text NOT NULL -- [[{"category_id": 1, "meal_id": 1}]], [day (0-6)][category/meal assignment]
);

CREATE TABLE meals (
  id                        INTEGER PRIMARY KEY,
  name                      text    NOT NULL,
  recipe_url                text    NOT NULL,
  preview_image_mode        text    NOT NULL DEFAULT "none",
  preview_image_url         text    NOT NULL DEFAULT "",
  preview_image_bytes       blob    NOT NULL DEFAULT "",
  ingredients_image_mode    text    NOT NULL DEFAULT "none",
  ingredients_image_url     text    NOT NULL DEFAULT "",
  ingredients_image_bytes   blob    NOT NULL DEFAULT "",
  ingredients               text    NOT NULL -- [ingredient ids]
);

CREATE TABLE ingredients (
  id INTEGER PRIMARY KEY,
  name text NOT NULL,
  ingredient_category_id INTEGER NOT NULL DEFAULT 0,
  want_override_count INTEGER NOT NULL DEFAULT 0,
  got_count INTEGER NOT NULL DEFAULT 0,
  shopped BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE ingredient_categories (
  id         INTEGER PRIMARY KEY,
  name       text    NOT NULL,
  sort_index INTEGER NOT NULL DEFAULT 0
)