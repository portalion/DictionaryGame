package room

import (
	"errors"
	"net/http"
	"server/internal/utils"

	"github.com/gorilla/websocket"
)

var	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {return true},
	}

type RoomHub struct {
	Rooms map[string]*RoomManager
}

func NewRoomHub() *RoomHub {
	return &RoomHub{
		Rooms: make(map[string]*RoomManager),
	}
}

func (rh *RoomHub) CreateRoom() string {
	const codeLength = 6
	code := utils.GenerateRandomCode(codeLength)

	for {
		if _, ok := rh.Rooms[code]; !ok {
			break
		}
		code = utils.GenerateRandomCode(codeLength)
	}

	rh.Rooms[code] = NewRoomManager()

	return code
}

func (rh *RoomHub) JoinRoom(w http.ResponseWriter, r *http.Request, code string, username string) error {
	if _, ok := rh.Rooms[code]; !ok {
		return errors.New("room with that code doesn't exist")
	}

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	data := WantsConnectionData{
		WsConnection: conn,
		Username: username,
	}

	rh.Rooms[code].Connect <- data
	return nil
}