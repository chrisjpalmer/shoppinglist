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
  want_override_count INTEGER NOT NULL DEFAULT 0
)