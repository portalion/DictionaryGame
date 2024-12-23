import { useEffect, useState } from "react";
import { hostname } from "../config";

async function CreateRoom(setCurrentRoomId: React.Dispatch<React.SetStateAction<string>>) {
    const response = await fetch(`http://${hostname}/room/create`, {method: "POST"})
    const data = await response.json() as string;
    setCurrentRoomId(data)
}

function RoomSelector(props: {setCurrentRoomId: React.Dispatch<React.SetStateAction<string>>, setUsername: React.Dispatch<React.SetStateAction<string>> })
{
    const [roomIdText, setRoomIdText] = useState("");
    const [username, setUsername] = useState("");

    return (<div>
        <div>
            <div>Paste code</div>
            <input type="text" maxLength={6} value={roomIdText} onChange={(e) => setRoomIdText(e.target.value)} placeholder="room code"/>
            <input type="text" maxLength={20} value={username} onChange={(e) => setUsername(e.target.value)} placeholder="username"/>
            <button onClick={() => {
                props.setCurrentRoomId(roomIdText)
                props.setUsername(username)
            }}>Join room</button>
        </div>
        <button onClick={() => CreateRoom(props.setCurrentRoomId)}>Create room</button>
    </div>);
}

export default RoomSelector;