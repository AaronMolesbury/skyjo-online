import { useWebSocket } from "./LobbyScreen"
import { useState } from "react";
import "./Button"
import Button from "./Button";

function ReadyUpButton() {
    const ws = useWebSocket();
    const [ready, setReady] = useState(false);

    if (!ws) {
        return <div>Loading...</div>
    }

    const readyUpClicked = () => {
        setReady(true)
        ws.send("ready");
    }

    let readyStateComponent = <Button clickHandler={readyUpClicked} labelText="Ready Up"/>;
    if (ready) {
        readyStateComponent = <div>Waiting for other players to ready up...</div>
    }

    return readyStateComponent
}

export default ReadyUpButton