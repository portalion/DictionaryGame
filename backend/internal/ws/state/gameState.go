package state

import (
	"server/internal/ws/event"
	"server/internal/ws/user"
)

type Player struct {
	Lives int
}

type GameState struct {
	Players map[*user.User] *Player
	usersOrder []*user.User
	currentPlayer *Player

	broadcast func(event.Event)
	swapState func(State)
}

func (gs* GameState) Start() {
	
}

func NewGameState(users []*user.User, broadcast func(event.Event), swapState func(State)) *GameState {
	gs := &GameState{broadcast: broadcast, swapState: swapState, usersOrder: users, Players: make(map[*user.User]*Player)}
	for _, user := range users {
		gs.Players[user] = &Player{Lives: 3}
	}

	if len(users) != 0 {
		gs.currentPlayer = gs.Players[users[0]]
	}

	return gs
}

func (gs* GameState) OnUserConnection(*user.User) {

} 

func (gs* GameState) OnUserDisconnection(*user.User) {

}

func (gs* GameState) ProcessMessage(message ServerStateMessage) error {
	return nil
}