const SET_PIXEL = 0
const SET_PALLETTE = 1

var pallette = []

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
    x = buffer[1];
    y = buffer[2];
    i = buffer[3];

    ctx.fillStyle = pallette[i]
    ctx.fillRect(x*3, y*3, 3, 3)
}

// buffer = type(1), index(1), r(1), g(1), b(1), ...
function setPallette(buffer) {
    for(i=1; i<buffer.byteLength; i+=4) {
        pallette[i] = hexCodeFromBytes(buffer[i+1], buffer[i+2], buffer[i+3]);
    }
}

function hexCodeFromBytes(r, g, b) {
    return "#" + r.toString(16) + g.toString(16) + b.toString(16);
}

init();
