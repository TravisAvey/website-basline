
// our container of all our categories buttons
const cats = document.getElementById("categories")
// our hidden input to put our values for hx-post to server
var catOutput = document.getElementById("categories-output")
// a list to keep track of all our categories
var categories = []

function updateCategories() {
  for (const cat of cats.children) {
    // get current category and index
    const current = cat.innerText
    const index = categories.indexOf(current)

    if (cat.classList.contains("badge-success")) {
      // if it's not already in list
      // add the category to list
      if (index == -1) {
       categories.push(current) 
      }
      catOutput.value = categories
    } 
  } 
}

for (const cat of cats.children) {
  cat.addEventListener("click", event => {

    // get current category and index
    const current = cat.innerText
    const index = categories.indexOf(current)

    // if the class is badge-outline, then it's not
    // a selected category
    if (cat.classList.contains("badge-outline")) {
      // flip the look of selected
      cat.classList.remove("badge-outline")
      cat.classList.add("badge-success")

      // if it's not already in list
      // add the category to list
      if (index == -1) {
       categories.push(current) 
      }

      // else its already a selected category
      // so flip the look..
    } else {
      cat.classList.remove("badge-success")
      cat.classList.add("badge-outline")

      // check if in (should be)
      if (index >= 0) {
        // remove the category
        categories.splice(index, 1)
      }
    }
    // regardless, update the list of categories
    // and put in our hidden value input
    catOutput.value = categories
  })
}
