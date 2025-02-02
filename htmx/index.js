
//Change text color on mouse over
// function mouseOver() {
//     var element = document.getElementById('tb-cell-icon-connected')
//     element.classList.replace("text-gray-400", "text-green-400")
//     animate()
//     console.log("animate")
// }


// Change text color on mouse leave
// function mouseLeave() {
//     var element = document.getElementById('tb-cell-icon-connected')
//     element.classList.replace("text-green-400", "text-gray-400")
// }

function mouseOver() {
    var element = document.getElementById('tb-cell-icon-connected')
    element.classList.replace("text-gray-400", "text-green-400")
}

function mouseLeave() {
    var element = document.getElementById('tb-cell-icon-connected')
    element.classList.replace("text-green-400", "text-gray-400")
}

//Animate connected SVG. Rotation and color
function animate(){
    const date = new Date()
    const second = date.getSeconds()
    const ico = document.getElementById('icon-connected')
    var secondLast
    var element = document.getElementById('tb-cell-icon-connected')

    ico.setAttribute('transform', `rotate(${(360/8)* second})`)
    requestAnimationFrame(animate) 

    if (secondLast != second){
        secondLast = second
        if (second % 2 == true){
            element.classList.replace("text-gray-400", "text-green-400")
        }else{
            element.classList.replace("text-green-400", "text-gray-400")
        }
    }
}

requestAnimationFrame(animate)
