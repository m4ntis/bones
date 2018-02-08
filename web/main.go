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
