package main

import (
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
	"encoding/json"
)

var active map[string]UserConnection

type Dialog []string

type UserConnection struct {
	name string
	conversations map[string]Dialog
}

func (c UserConnection) SendMessage(author string, message string) {
	
}

func main() {
	active = make(map[string]UserConnection)

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		currentConnect := UserConnection{}

		so.Join("chat")

		so.On("auth", func(name string) {
			currentConnect.name = name
			active[name] = currentConnect

			user_names := make([]string, len(active))
			i := 0
			for _, v := range active {
				user_names[i] = v.name
				i++
			}
			log.Printf("%q\n", user_names)
			users_json, err := json.Marshal(user_names)
			if err != nil {
				log.Print(err)
				return
			}

			so.Emit("auth_success", string(users_json))
			so.BroadcastTo("chat", "new_user", name)
		})

		so.On("disconnection", func() {
			log.Println("on disconnect")
			so.BroadcastTo("chat", "user_disconnect", currentConnect.name)
			delete(active, currentConnect.name)
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