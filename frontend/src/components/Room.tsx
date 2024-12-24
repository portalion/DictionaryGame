import { useEffect, useMemo, useState } from 'react';
import { hostname } from '../config';
import useWebSocket from 'react-use-websocket';

class Event {
  type: string = '';
  payload: unknown;
}

function Room(props: {
  currentRoomId: string;
  setCurrentRoomId: React.Dispatch<React.SetStateAction<string>>;
  username: string;
}) {
  const wsUrl = `ws://${hostname}/ws/room/${props.currentRoomId}/join?username=${props.username}`;
  const { lastJsonMessage, sendJsonMessage } = useWebSocket(wsUrl, {
    share: true,
    onOpen: () => {
      console.log(`connect`);
      const message = new Event();
      message.type = 'room_state_requested';
      sendJsonMessage(message);
    },
    onClose: () => {
      props.setCurrentRoomId('');
      console.log(`disconnect`);
    },
    onError: () => props.setCurrentRoomId(''),
  });
  const [users, setUsers] = useState<string[]>([props.username]);

  useEffect(() => {
    if (lastJsonMessage !== null) {
      const message = lastJsonMessage as Event;
      console.log(message);
      if (message.type === 'user_joined')
        setUsers(users.concat((message.payload as { username: string }).username));
      else if (message.type === 'user_disconnected')
        setUsers(users.filter(v => v != (message.payload as { username: string }).username));
      else if (message.type === 'room_state')
        setUsers((message.payload as { users: string[] }).users);
    }
  }, [lastJsonMessage]);

  return (
    <div>
      Code: {props.currentRoomId}
      <button
        onClick={() => {
          props.setCurrentRoomId('');
        }}>
        Disconnect
      </button>
      {users.map(v => (
        <div>{v}</div>
      ))}
    </div>
  );
}

export default Room;
