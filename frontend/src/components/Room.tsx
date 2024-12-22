import { useEffect, useMemo, useState } from "react";
import { hostname } from "../config";
import useWebSocket from "react-use-websocket";

class Event
{
    type: string = ""
    payload: unknown = ""
}

function Room(props: {currentRoomId: string, setCurrentRoomId: React.Dispatch<React.SetStateAction<string>>})
{
    const wsUrl = `ws://${hostname}/ws/room/join/${props.currentRoomId}`
    const { lastJsonMessage } = useWebSocket(wsUrl, { share: true,  
        onOpen: () => console.log(`connect `),
		onClose: () => console.log(`disconnect `) });
    const [users, setUsers] = useState<boolean[]>([true]);

    useEffect(() =>
    {
        if (lastJsonMessage !== null)
        {
            const message = lastJsonMessage as Event;
            if (message.type === 'user_joined')
                setUsers(users.concat(true));
            else if (message.type === 'user_disconnected')
                setUsers(users.slice(1))
        }
    }, [lastJsonMessage])

    return (
    <div>
        <button onClick={() => {
            props.setCurrentRoomId("")
            }}>Disconnect</button>
        {users.map(_ => <div>user</div>)}
    </div>);
}

export default Room;