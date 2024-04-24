
function showDeleteDialog(id) {
  const ele = "modal_" + id
  document.getElementById(ele).showModal()
}

function closeDeleteDialog(id) {
  const ele = "modal_" + id
  document.getElementById(ele).close()
}

function removeElement(ele) {
  document.getElementById(ele).remove()
}

window.onload = function () {
  document.body.addEventListener("htmx:confirm", function(e) {
    //console.log("confirm button pressed")
    //console.log(e.detail.path)
  })
  
}
