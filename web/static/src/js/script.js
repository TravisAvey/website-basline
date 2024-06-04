
function showDeleteDialog(id) {
  const ele = "modal_" + id
  document.getElementById(ele).showModal()
}

function closeDeleteDialog(id) {
  const ele = "modal_" + id
  document.getElementById(ele).close()
}

function removeCardElement(id) {
  const ele = "card-" + id
  document.getElementById(ele).remove()
}

function removeMsgElement(id) {
  const ele = "message-" + id
  document.getElementById(ele).remove()
}

function showMessageModal(id) {
  const ele = "msg-modal-" + id
  document.getElementById(ele).showModal()
}

function closeMessageModal(id) {
  const ele = "msg-modal-" + id
  document.getElementById(ele).close()
}

window.onload = function () {
  document.body.addEventListener("htmx:confirm", function(e) {
    //console.log("confirm button pressed")
    //console.log(e.detail.path)
  })
}


function copyToClipboard(id) {

  const ele = "imgurl-" + id
  var text = document.getElementById(ele)

  navigator.clipboard.writeText(text.innerText)

  var classes = text.previousElementSibling.classList
  classes.add("tooltip-open", "transition", "transition-opacity", "duration-500", "ease-in-out")

  setTimeout(() => {
    classes.remove("tooltip-open", "transition", "transition-opacity", "duration-500", "ease-in-out") 
  }, 1000);
}

