import Button from "./Button";
import { useData, useWebSocket } from "./LobbyScreen";
import "../css/EndGameScreen.css"

function EndGameScreen() {
    const ws = useWebSocket();
    const data = useData();

    const resetGameClicked = () => {
        ws?.send("reset");
    }

    return (
        <div className="EndGameScreen">
            <div className="EndGameScreen_WinnerText">Player {data?.winnerId} wins!</div>
            {data?.playerId === 0 ? (
                <Button clickHandler={resetGameClicked} labelText={"Restart?"}></Button>
            ) : (
                <div className="EndGameScreen_WaitingText">Waiting for the host to start a new game...</div>
            )}
        </div>
    )
}

export default EndGameScreen