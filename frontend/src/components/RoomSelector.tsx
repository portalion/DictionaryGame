import { useState } from "react";

function RoomSelector(props: {setCurrentRoomId: React.Dispatch<React.SetStateAction<string>>})
{
    const [roomIdText, setRoomIdText] = useState("");

    return (<div>
        <div>
            <div>Paste code</div>
            <input type="text" maxLength={6} value={roomIdText} onChange={(e) => setRoomIdText(e.target.value)}/>
            <button onClick={() => {
                props.setCurrentRoomId(roomIdText)
            }}>Join room</button>
        </div>
        <button>Create room</button>
    </div>);
}

export default RoomSelector;