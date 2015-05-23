package main

import (
	"log"
	"net/http"
	"github.com/googollee/go-socket.io"
	"encoding/json"
	"errors"
	"strconv"
)

func getContactList(userId int) (data contactList, err error) {
	data, ok := users[userId]
	if !ok {
		return data, errors.New("No such user")
	}
	return data, nil
}

func updateStatus(so *socketio.Socket, data contactData) error {
	dataMsg, err := packEvent(data)
	if err != nil {
		return err
	}

	err = (*so).Emit(CONTACT_STATUS, string(dataMsg))
	if err != nil {
		return err
	}
	return nil
}

const AUTH_REQUEST = "auth_request"

type authRequest struct {
	Ssid string
}

const AUTH_RESPONSE = "auth_response"

type contactList []contactData

type contactData struct {
	Name string
	UserId int
	Online bool
}

type authResponse struct {
	Result bool
	UserId int
	Contacts contactList
}

func checkAuth(ssid string) (ok bool, userId int, err error) {
	userId, err = strconv.Atoi(ssid)
	if err != nil {
		return false, 0, err
	}
	if userId < 100 || userId > 104 {
		return false, 0, nil
	}
	return true, userId, nil
}

const CONTACT_STATUS = "contact_status"
const MESSAGE_SEND = "message_send"
const MESSAGE_RECV = "message_recv"
const MESSAGE_ACCEPTED = "message_accepted"
const MESSAGE_READ = "message_read"
const NEW_CONTACT = "new_contact"
const DELETE_CONTACT = "delete_contact"

func unpackEvent(msg string, p interface{}) (err error) {
	return json.Unmarshal([]byte(msg), &p)
}

func packEvent(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

var users map[int]contactList

func main() {

	users = make(map[int]contactList)

	users[101] = contactList{
		contactData{"ivan", 102, true},
		contactData{"petr", 103, true},
		contactData{"vladimir", 104, true},
	}

	users[102] = contactList{
		contactData{"huilo", 101, true},
		contactData{"petr", 103, true},
		contactData{"vladimir", 104, true},
	}

	users[103] = contactList{
		contactData{"ivan", 102, true},
		contactData{"huilo", 101, true},
		contactData{"vladimir", 104, true},
	}

	users[104] = contactList{
		contactData{"ivan", 102, true},
		contactData{"petr", 103, true},
		contactData{"huilo", 101, true},
	}


	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {

		so.Join("chat")

		so.On(AUTH_REQUEST, func(msg string) {

			var request authRequest

			err := unpackEvent(msg, &request)

			if err != nil {
				log.Print(err)
				return
			}

			ok, userId, err := checkAuth(request.Ssid)
			if err != nil {
				log.Print(err)
				return
			}
			if !ok {
				respMsg, err := packEvent(authResponse{false, 0, nil})
				if err != nil {
					log.Print(err)
					return
				}
				err = so.Emit(AUTH_RESPONSE, string(respMsg))
				if err != nil {
					log.Print(err)
					return
				}
			}

			list, err := getContactList(userId)
			if err != nil {
				log.Print(err)
				return
			}
			respMsg, err := packEvent(authResponse{true, userId, list})
			if err != nil {
				log.Print(err)
				return
			}
			err = so.Emit(AUTH_RESPONSE, string(respMsg))
			if err != nil {
				log.Print(err)
				return
			}

			err = updateStatus(&so, list[0])
			if err != nil {
				log.Print(err)
				return
			}
		})

		so.On("disconnection", func() {
			log.Printf("%q\n", so)
			log.Println("on disconnect")
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