
function showDeleteDialog(id) {
  const ele = "modal_" + id
  document.getElementById(ele).showModal()
}

function closeDeleteDialog(id) {
  const ele = "modal_" + id
  document.getElementById(ele).close()
}

function removePostElement(id) {
  const ele = "post-" + id
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
  document.body.addEventListener("htmx:afterRequest", function (evt) {

    if (evt.detail.xhr.status == 401) {
      const ele = document.getElementById("response")
      ele.attributes[0].textContent = "{ open: true }"
      const msgEle = document.getElementById("message")
      res = JSON.parse(evt.detail.xhr.response)
      msgEle.innerHTML = res.content

    }
  })  
}

