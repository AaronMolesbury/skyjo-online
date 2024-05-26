import { useData, useWebSocket } from "../App"
import "../css/GameWindow.css"

function GameWindow() {
    const ws = useWebSocket();
    const data = useData();

    if (!ws) {
        return (
            <div>Loading...</div>
        )
    }

    const readyUpClicked = () => {
        console.log("clicked deez")
        ws.send("ready");
    }

    if (!data) {
        return (
            <div className="GameWindow">
                <div className="GameWindow_DeckWrapper">
                </div>
                <div className="GameWindow_ButtonWrapper">
                    <div className="GameWindow_Button" onClick={readyUpClicked}>
                        <div className="GameWindow_ButtonLabel">
                            Ready Up
                        </div>
                    </div>
                </div>
            </div>
        )
    }

    const playerId = data.playerId;
    const hand = data.players[playerId].hand;
   
    const cardsComponent = hand.map((row, i) => (
        <div className="GameWindow_CardRow" key={i}>
            {row.map((card, i) => (
                <div className="GameWindow_Card" key={i}>
                    {card?.value}
                </div>
            ))}
        </div>
    ))

    return (
        <div className="GameWindow">
            <div className="GameWindow_DeckWrapper">
                {cardsComponent}
            </div>
            <div className="GameWindow_ButtonWrapper">
                <div className="GameWindow_Button" onClick={readyUpClicked}>
                    <div className="GameWindow_ButtonLabel">
                        Ready Up
                    </div>
                </div>
            </div>
        </div>
    )
}

export default GameWindow