import { useState } from 'react';
import { hostname } from '../config';
import { GlobalData, GlobalStateAction } from '../App.types';

type RoomSelectorProps = {
  dispatch: React.Dispatch<GlobalStateAction>;
  globalState: GlobalData;
};

async function CreateRoom(): Promise<string> {
  const response = await fetch(`http://${hostname}/room/create`, { method: 'POST' });
  const data = (await response.json()) as { code: string };
  return data.code;
}

function RoomSelector({ dispatch, globalState }: RoomSelectorProps) {
  const [roomIdText, setRoomIdText] = useState('');
  const [username, setUsername] = useState(globalState.username);

  return (
    <div>
      <div>
        <div>Paste code</div>
        <input
          type="text"
          maxLength={6}
          value={roomIdText}
          onChange={e => setRoomIdText(e.target.value)}
          placeholder="room code"
        />
        <input
          type="text"
          maxLength={20}
          value={username}
          onChange={e => setUsername(e.target.value)}
          placeholder="username"
        />
        <button
          onClick={() => {
            dispatch({
              type: 'ChangeGlobalData',
              newData: { username: username, currentRoomId: roomIdText },
            });
          }}>
          Join room
        </button>
      </div>
      <button
        onClick={async () => {
          dispatch({
            type: 'ChangeGlobalData',
            newData: { username: username, currentRoomId: await CreateRoom() },
          });
        }}>
        Create room
      </button>
    </div>
  );
}

export default RoomSelector;
