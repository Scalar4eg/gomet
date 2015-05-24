package main

import (
	"github.com/googollee/go-socket.io"
	"log"
)


type Message struct {
	Text string
	UserIdFrom int
	UserIdTo int
	MessageId int
}

type MessageId struct {
	Value int
}

var awaitingReadMessages map[int]int

const MESSAGE_SEND = "message_send" 		//from ClientA to server
const MESSAGE_ACCEPTED = "message_accepted" //from server to ClientA with message_id
const MESSAGE_RECV = "message_recv" 		//from server to ClientB
const MESSAGE_READ = "message_read" 		//from ClientB to server and from Server to ClientA

func init () {
	awaitingReadMessages = make(map[int]int)
}

func saveMessage(message Message) (int, error) {
	mockMessageId++
	 return mockMessageId, nil
}

func markMessageAsRead(messageId int) error {
	return nil
}

func onSocketMessageRead(msg string, so socketio.Socket) {
	var messageId MessageId
	err := unpackEvent(msg, &messageId)

	if err != nil {
		log.Print(err)
		return
	}

	err = markMessageAsRead(messageId.Value)


	if err != nil {
		log.Print(err)
		return
	}


	userIdAwaiting, ok := awaitingReadMessages[messageId.Value]

	if !ok {
		return
	}

	so, ok = getSocketOfUser(userIdAwaiting)

	if !ok {
		return
	}

	readResponse, err := packEvent(messageId)

	if err != nil {
		log.Print(err)
		return
	}
	so.Emit(MESSAGE_READ, readResponse)
}

func onSocketMessageSend(userId int, msg string, so socketio.Socket) {
	var message Message
	err := unpackEvent(msg, &message)

	if err != nil {
		log.Print(err)
		return
	}

	message.UserIdFrom = userId

	messageId, err := saveMessage(message)

	if err != nil {
		log.Print(err)
		return
	}

	message.MessageId = messageId

	acceptResponse, err := packEvent(MessageId{message.MessageId})

	if err != nil {
		log.Print(err)
		return
	}

	err = so.Emit(MESSAGE_ACCEPTED, acceptResponse)

	if err != nil {
		log.Print(err)
		return
	}

	recvSocket, ok := getSocketOfUser(message.UserIdTo)

	if !ok {
		log.Print("user is offline")
		return
	}

	recvResponse, err := packEvent(message)


	if err != nil {
		log.Print(err)
		return
	}

	err = recvSocket.Emit(MESSAGE_RECV, recvResponse)

	if err != nil {
		log.Print(err)
		return
	}

	awaitingReadMessages[messageId] = userId
}



