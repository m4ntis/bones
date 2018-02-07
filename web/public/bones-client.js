const SET_PIXEL = 0
const SET_PALLETTE = 1

function init_canvas(name) {
    var canvas = document.getElementById(name);
    var ctx = canvas.getContext("2d");

    ctx.fillStyle = "#FF0000";
    ctx.fillRect(0,0,768,672);

    return [canvas, ctx];
}

function init() {
    [canvas, ctx] = init_canvas("screen");
    
    con = new WebSocket("ws://" + document.location.host + "/ws");
    con.onmessage = onMessage;
}

function onMessage(evt) {
    var buffer = new Uint8Array(evt.data);
    if buffer[0] == SET_PIXEL {
        setPixel(buffer);
    else {
        setPallette(buffer);
    }
}

// buffer = type(1), x(1), y(1), colour(1)
function setPixel(buffer) {
}

function setPallette(buffer) {
}

init();
