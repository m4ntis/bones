package main

import (
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
	b := []byte{1, 0, 255, 0, 0, 1, 0, 255, 0, 2, 0, 0, 255}
	return b
}

func getColours() [][]byte {
	return [][]byte{
		[]byte{0, 0, 0, 0},
		[]byte{0, 0, 1, 0},
		[]byte{0, 0, 2, 0},
		[]byte{0, 0, 5, 0},
		[]byte{0, 0, 7, 0},
		[]byte{0, 0, 10, 0},
		[]byte{0, 0, 14, 0},

		[]byte{0, 1, 1, 0},
		[]byte{0, 1, 4, 0},
		[]byte{0, 1, 5, 0},
		[]byte{0, 1, 6, 0},
		[]byte{0, 1, 7, 0},
		[]byte{0, 1, 8, 0},
		[]byte{0, 1, 10, 0},
		[]byte{0, 1, 11, 0},
		[]byte{0, 1, 13, 0},
		[]byte{0, 1, 14, 0},

		[]byte{0, 2, 1, 0},
		[]byte{0, 2, 4, 0},
		[]byte{0, 2, 5, 0},
		[]byte{0, 2, 6, 0},
		[]byte{0, 2, 7, 0},
		[]byte{0, 2, 8, 0},
		[]byte{0, 2, 10, 0},
		[]byte{0, 2, 12, 0},
		[]byte{0, 2, 14, 0},

		[]byte{0, 3, 1, 0},
		[]byte{0, 3, 5, 0},
		[]byte{0, 3, 6, 0},
		[]byte{0, 3, 7, 0},
		[]byte{0, 3, 10, 0},
		[]byte{0, 3, 14, 0},

		[]byte{0, 4, 0, 0},
		[]byte{0, 4, 1, 0},
		[]byte{0, 4, 2, 0},
		[]byte{0, 4, 6, 0},
		[]byte{0, 4, 10, 0},
		[]byte{0, 4, 14, 0},
	}
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

func handleClient(ws *websocket.Conn) {
	time.Sleep(100 * time.Millisecond)
	ws.WriteMessage(websocket.BinaryMessage, getPallette())

	for _, colour := range getColours() {
		ws.WriteMessage(websocket.BinaryMessage, colour)
		time.Sleep(200 * time.Millisecond)
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

	go handleClient(ws)
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
