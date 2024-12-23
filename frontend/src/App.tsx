import { useState } from 'react';
import RoomSelector from './components/RoomSelector';
import Room from './components/Room';

function App() {
  const [currentRoomId, setCurrentRoomId] = useState('');
  const [username, setUsername] = useState('');

  return (
    <>
      {!currentRoomId ? (
        <RoomSelector
          setCurrentRoomId={setCurrentRoomId}
          setUsername={setUsername}
          username={username}
        />
      ) : (
        <Room
          currentRoomId={currentRoomId}
          setCurrentRoomId={setCurrentRoomId}
          username={username}
        />
      )}
    </>
  );
}

export default App;
