package main
import (
	"strconv"
	"github.com/googollee/go-socket.io"
)

const AUTH_REQUEST = "auth_request"

type authRequest struct {
	Ssid string
}

const AUTH_RESPONSE = "auth_response"


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

func onSocketAuthRequest (msg string, so socketio.Socket) (int, error) {

	var request authRequest

	err := unpackEvent(msg, &request)

	if err != nil {
		return 0, err
	}

	ok, userId, err := checkAuth(request.Ssid)
	if err != nil {
		return 0, err
	}
	if !ok {
		respMsg, err := packEvent(authResponse{false, 0, nil})
		if err != nil {
			return 0, err
		}
		err = so.Emit(AUTH_RESPONSE, respMsg)
		if err != nil {
			return 0, err
		}
	}



	list, err := getContactList(userId)
	if err != nil {
		return 0, err
	}
	respMsg, err := packEvent(authResponse{true, userId, list})
	if err != nil {
		return 0, err
	}
	err = so.Emit(AUTH_RESPONSE, respMsg)
	if err != nil {
		return 0, err
	}

	err = updateStatus(&so, list[0])
	if err != nil {
		return 0, err
	}

	err = deleteContact(&so, 102)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
