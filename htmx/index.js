
//Change text color on mouse over
function mouseOver() {
    var element = document.getElementById('tb-cell-icon-connected')
    element.classList.replace("text-gray-400", "text-green-400")
}

//Change text color on mouse leave
function mouseLeave() {
    var element = document.getElementById('tb-cell-icon-connected')
    element.classList.replace("text-green-400", "text-gray-400")
}