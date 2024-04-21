
function showDeleteDialog(id) {
  const ele = "modal_" + id
  document.getElementById(ele).showModal()
}

function closeDeleteDialog(id) {
  const ele = "modal_" + id
  document.getElementById(ele).close()
}

function removeElement(id) {
  const ele = "post-" + id
  document.getElementById(ele).remove()
}
