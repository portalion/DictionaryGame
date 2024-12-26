import { useEffect, useState } from 'react';
import { hostname } from '../config';
import useWebSocket from 'react-use-websocket';
import { GlobalData, GlobalStateAction } from '../App.types';

class Event {
  type: string = '';
  payload: unknown;
}

type RoomProps = {
  dispatch: React.Dispatch<GlobalStateAction>;
  globalState: GlobalData;
};

function Room({ dispatch, globalState }: RoomProps) {
  const wsUrl = `ws://${hostname}/ws/room/${globalState.currentRoomId}/join?username=${globalState.username}`;
  const { lastJsonMessage, sendJsonMessage } = useWebSocket(wsUrl, {
    share: true,
    onOpen: () => {
      console.log(`connect`);
      const message = new Event();
      message.type = 'room_state_requested';
      sendJsonMessage(message);
    },
    onClose: () => {
      dispatch({ type: 'ChangeRoomCode', newRoomCode: '' });
      console.log(`disconnect`);
    },
    onError: () => dispatch({ type: 'ChangeRoomCode', newRoomCode: '' }),
  });
  const [users, setUsers] = useState<string[]>([globalState.username]);

  useEffect(() => {
    if (lastJsonMessage !== null) {
      const message = lastJsonMessage as Event;
      switch (message.type) {
        case 'user_joined': {
          setUsers(users.concat((message.payload as { username: string }).username));
          break;
        }
        case 'user_disconnected': {
          setUsers(users.filter(v => v != (message.payload as { username: string }).username));
          break;
        }
        case 'room_state': {
          setUsers((message.payload as { users: string[] }).users);
          break;
        }
        case 'game_started': {
          break;
        }
      }
    }
  }, [lastJsonMessage]);

  return (
    <div>
      Code: {globalState.currentRoomId}
      <button
        onClick={() => {
          dispatch({ type: 'ChangeRoomCode', newRoomCode: '' });
        }}>
        Disconnect
      </button>
      {users.map(v => (
        <div>{v}</div>
      ))}
      <button
        onClick={() => {
          const message = new Event();
          message.type = 'game_start';
          sendJsonMessage(message);
        }}>
        Start Game
      </button>
    </div>
  );
}

export default Room;
