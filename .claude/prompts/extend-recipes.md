Work on the frontend project to extend the recipes page so that the ingredients image can be edited on the page. The recipes page is located in the `frontend` directory in `src/routes/recipes/[meal]/+page.svelte`. Adapt code from `src/routes/meals/+page.svelte` to do this. The backend accepts two requests to save a meal:
  1. Save the meal  (includes metadata about the image - whether it is internal/external)
  2. Save the meal image
  
It should be noted that the recipe page should not only support setting the image to a file that the user selects (internal image), but should also support inputting a URL (external image). Allow the user to choose either of these options by providing two distinct buttons on the recipe page (locate the buttons in the top left corner of the ingredients image). For now use the `<Button>` class to do this, avoid icons, text is fine. Use a flex box with flex column so that the buttons stack vertically.

If the user selects the "Upload Image" button this is the flow:
1. The user is presented with a file dialog
2. The user selects the file
3. A request to save the meal is sent to the backend - saves the entire meal but sets the image to internal + sets the url 
4. A request to upload the image is sent to the backend

If the user selects "Set image url" button:
1. The user is presented with a modal with a textbox
2. The user enters the image url
3. A request to save the meal is sent to the backend - saves the entire meal but sets the image to external + sets the url