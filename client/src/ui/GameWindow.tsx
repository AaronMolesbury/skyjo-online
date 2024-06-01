import { useData, useWebSocket } from "./LobbyScreen"
import "../css/GameWindow.css"
import Button from "./Button";
import CardGrid from "./CardGrid";
import ScoreLabel from "./ScoreLabel";
//import Telemetry from "./Telemetry";
import { useEffect, useState } from "react";
import Piles from "./Piles";
import CardInHand from "./CardInHand";

function GameWindow() {
    const ws = useWebSocket();
    const data = useData();

    const [width, setWidth] = useState(window.innerWidth);

    useEffect(() => {
        const handleResize = () => {
            setWidth(window.innerWidth);
        }

        window.addEventListener("resize", handleResize);
        return () => {
            window.removeEventListener("resize", handleResize);
        }
    }, []);

    if (!ws) {
        // TEMP while no landing screen
        return (
            <div>Loading...</div>
        )
    }

    const resetGameClicked = () => {
        ws.send("reset");
    }

    if (data?.winnerId !== -1) {
        // TEMP while no landing screen
        return (
            <>
                <div>Player {data?.winnerId} wins!</div>
                <Button clickHandler={resetGameClicked} labelText={"Restart?"}></Button>
            </>
        )
    }

    const playerId = data.playerId;
    const player = data.players[playerId]
    
    let deckWrapper: JSX.Element;
    if (width <= 1250) {
        deckWrapper = (
            <>
                <div className="GameWindow_DeckWrapper">
                    <CardGrid hand={player.hand}/>
                    <div className="GameWindow_PileWrapper">
                        <Piles faceUpCard={data.lastDiscardedCard}/>
                        <CardInHand faceUpCard={data.cardInHand}/>
                    </div>
                    <ScoreLabel/>
                </div>
            </>
        )
    } else {
        deckWrapper = (
            <>
                <div className="GameWindow_DeckWrapper">
                    <Piles faceUpCard={data.lastDiscardedCard}/>
                    <CardGrid hand={player.hand}/>
                    <CardInHand faceUpCard={data.cardInHand}/>
                    <ScoreLabel/>
                </div>
            </>
        )
    }

    return (
        <div className="GameWindow">
            {/* <Telemetry 
                playerServerId={playerId}
                gameState={player.turnType}
                cardInHandValue={data.cardInHand?.value ?? null}
                lastDiscardedCardValue={data.lastDiscardedCard?.value ?? null}
            /> */}
            {deckWrapper}
        </div>
    )
}

export default GameWindow