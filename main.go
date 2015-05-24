package main

import (
	"log"
	"net/http"
	"github.com/googollee/go-socket.io"
	"encoding/json"
)

func unpackEvent(msg string, p interface{}) (err error) {
	return json.Unmarshal([]byte(msg), &p)
}

func packEvent(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

var userConnections map[int]socketio.Socket

func getSocketOfUser (userId int) (socketio.Socket, bool) {
	so, ok := userConnections[userId]
	return so, ok
}

func main() {

	userConnections = make(map[int]socketio.Socket)
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {

		var currentUserId int
		so.Join("chat")

		so.On(AUTH_REQUEST, func (msg string) {
				currentUserId, err = onSocketAuthRequest(msg, so)
				if err != nil {
					log.Print(err)
					return
				}
				userConnections[currentUserId] = so
			})

		so.On(MESSAGE_SEND, func (msg string) {
				onSocketMessageSend(currentUserId, msg, so)
			})


		so.On(MESSAGE_READ, func (msg string) {
				onSocketMessageRead(msg, so)
			})
		so.On("disconnection", func() {
			log.Printf("%q\n", so)
			log.Println("on disconnect")
			if currentUserId !=0 {
				delete(userConnections, currentUserId)
			}

		})

		so.On("chat_message", func(msg string) {

		})

	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
