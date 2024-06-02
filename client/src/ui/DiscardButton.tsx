import { useWebSocket } from "./LobbyScreen";
import { usePlayer } from "../util/hooks";
import "../css/DiscardButton.css";

function DiscardButton() {
    const ws = useWebSocket();
    const {player, isCurrentTurn} = usePlayer();

    const discardClicked = () => {
        if (isCurrentTurn) {
            ws?.send("discard");
        }
    }

    return (
        <>
            {player?.turnType === "swap-discard" ? <div className="DiscardButton" onClick={discardClicked}/> : <></>}
        </>
    );
}

export default DiscardButton