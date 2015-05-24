package main

import (
	"github.com/googollee/go-socket.io"
	"errors"
)

const CONTACT_STATUS = "contact_status"
const CONTACT_DELETE = "delete_contact"

type contactList []contactData

type contactData struct {
	Name   string
	UserId int
	Online bool
}

type contactDeleteRequest struct {
	UserId int
}

func getContactList(userId int) (data contactList, err error) {
	data, ok := users[userId]
	if !ok {
		return data, errors.New("No such user")
	}
	return data, nil
}

func deleteContact(so *socketio.Socket, contactId int) error {
	request := contactDeleteRequest{contactId}
	data, err := packEvent(request)
	if err != nil {
		return err
	}
	return (*so).Emit(CONTACT_DELETE, data)
}

func updateStatus(so *socketio.Socket, data contactData) error {
	dataMsg, err := packEvent(data)
	if err != nil {
		return err
	}

	err = (*so).Emit(CONTACT_STATUS, dataMsg)
	if err != nil {
		return err
	}
	return nil
}
