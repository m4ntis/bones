package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

var (
	log      = logrus.WithField("cmd", "web")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func getPallette() []byte {
	// set pallette command, [colour num, r, g, b]...
	b := []byte{1, 0, 255, 0, 0, 1, 0, 255, 0, 2, 0, 0, 255}
	return b
}

func getColour() []byte {
	x := byte(rand.Intn(256))
	y := byte(rand.Intn(256))
	i := byte(rand.Intn(3))

	colour := make([]byte, 4)
	colour[0] = 0
	colour[1] = x
	colour[2] = y
	colour[3] = i

	return colour
}

func sendPallette(pallette map[color.NRGBA]bool, ws *websocket.Conn) map[color.NRGBA]int {
	internalPallette := make(map[color.NRGBA]int)
	data := []byte{1}

	i := 0
	for color := range pallette {
		data = append(data, byte(i))
		data = append(data, color.R)
		data = append(data, color.G)
		data = append(data, color.B)

		internalPallette[color] = i
		i++
	}

	ws.WriteMessage(websocket.BinaryMessage, data)
	return internalPallette
}

func sendImage(img image.Image, pallette map[color.NRGBA]int, ws *websocket.Conn) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.At(x, y)
			data := make([]byte, 4)
			data[0] = 0
			data[1] = byte(x)
			data[2] = byte(y)
			data[3] = byte(pallette[c.(color.NRGBA)])
			fmt.Println(data)
			ws.WriteMessage(websocket.BinaryMessage, data)
			time.Sleep(250 * time.Microsecond)
		}
	}
}

func handleImage(r io.Reader, ws *websocket.Conn) {
	pallette := make(map[color.NRGBA]bool)

	img, _ := png.Decode(r)
	b := img.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.At(x, y)

			if !pallette[c.(color.NRGBA)] {
				pallette[c.(color.NRGBA)] = true
			}
		}
	}

	internalPallette := sendPallette(pallette, ws)
	sendImage(img, internalPallette, ws)
}

func handleClientSlideshow(ws *websocket.Conn) {
	time.Sleep(100 * time.Millisecond)
	files, _ := ioutil.ReadDir("./public/slides/")
	for _, file := range files {
		f, _ := os.Open(fmt.Sprintf("./public/slides/%s", file.Name()))
		handleImage(f, ws)
		f.Close()
	}
}

func handleClient(ws *websocket.Conn) {
	time.Sleep(100 * time.Millisecond)
	ws.WriteMessage(websocket.BinaryMessage, getPallette())

	for {
		ws.WriteMessage(websocket.BinaryMessage, getColour())
		time.Sleep(2 * time.Millisecond)
	}
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		m := "Unable to upgrade to websockets"
		log.WithField("err", err).Println(m)
		http.Error(w, m, http.StatusBadRequest)
		return
	}

	go handleClientSlideshow(ws)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.WithField("PORT", port).Fatal("$PORT must be set")
	}

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", handleWebsocket)
	log.Println(http.ListenAndServe(":"+port, nil))
}
