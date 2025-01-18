

function mouseOver() {
    var element = document.getElementById('tb-cell-icon-connected')
    element.classList.replace("text-gray-400", "text-green-400")
}

function mouseLeave() {
    var element = document.getElementById('tb-cell-icon-connected')
    element.classList.replace("text-green-400", "text-gray-400")
}