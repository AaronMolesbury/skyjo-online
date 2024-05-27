import { useWebSocket } from "../App"
import { useState } from "react";
import "../css/LobbyWindow.css"
import "./Button"
import Button from "./Button";

function LobbyWindow() {
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
        readyStateComponent = <div>waiting for players</div>
    }

    return (
        <div className="LobbyWindow">
            <div className="LobbyWindow_Title">TUMGA</div>
            {readyStateComponent}
        </div>
    )
}

export default LobbyWindow