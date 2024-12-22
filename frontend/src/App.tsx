import { useState } from "react";
import RoomSelector from "./components/RoomSelector";
import Room from "./components/Room";

function App()
{
  const [currentRoomId, setCurrentRoomId] = useState<string>("");

  return (
    <>
        { !currentRoomId 
          ? <RoomSelector setCurrentRoomId={setCurrentRoomId}/> 
          : <Room currentRoomId={currentRoomId} setCurrentRoomId={setCurrentRoomId}/>
        }
    </>
  );

}

export default App
