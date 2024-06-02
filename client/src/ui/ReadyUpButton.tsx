import { useState } from "react";
import { useWebSocket } from "./LobbyScreen"
import "./Button"
import Button from "./Button";

function ReadyUpButton() {
    const ws = useWebSocket();
    const [ready, setReady] = useState(false);

    const readyUpClicked = () => {
        setReady(true);
        ws?.send("ready");
    }

    return (
        <>
            {ready ? <div>Waiting for other players to ready up...</div> : <Button clickHandler={readyUpClicked} labelText="Ready Up"/>}
        </>
    );
}

export default ReadyUpButton