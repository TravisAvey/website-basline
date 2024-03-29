
function showDeleteDialog(id) {
  ele = "modal_" + id
  document.getElementById(ele).showModal()
}

function closeDeleteDialog(id) {
  ele = "modal_" + id
  document.getElementById(ele).close()
}
