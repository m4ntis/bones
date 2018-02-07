function init_canvas(name) {
    var canvas = document.getElementById(name);
    var ctx = canvas.getContext("2d");

    ctx.fillStyle = "#FF0000";
    ctx.fillRect(0,0,768,672);

    return [canvas, ctx];
}

function init() {
    [canvas, ctx] = init_canvas("screen");
}

init();
