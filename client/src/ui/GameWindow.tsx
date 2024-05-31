import { useData, useWebSocket } from "../App"
import "../css/GameWindow.css"
import Button from "./Button";
import CardGrid from "./CardGrid";
import Pile from "./Pile";
import ScoreLabel from "./ScoreLabel";
import Telemetry from "./Telemetry";
import ForceSwap from "./gamestates/ForceSwap";
import SwapDiscard from "./gamestates/SwapDiscard";
import TakeFrom from "./gamestates/TakeFrom";

function GameWindow() {
    const ws = useWebSocket();
    const data = useData();

    if (!ws) {
        // TEMP while no landing screen
        return (
            <div>Loading...</div>
        )
    }

    const readyUpClicked = () => {
        ws.send("ready");
    }

    if (!data) {
        return (
            <div className="GameWindow">
                <div className="GameWindow_ButtonWrapper">
                   <Button clickHandler={readyUpClicked} labelText="Ready Up"/>
                </div>
            </div>
        )
    }

    if (data.winnerId !== -1) {
        // TEMP while no landing screen
        return (
            <div>Player {data.winnerId} wins!</div>
        )
    }

    const playerId = data.playerId;
    const player = data.players[playerId]
    
    let currentGameState: JSX.Element;
    switch (player.turnType) {
        case "take-from":
            currentGameState = (
                <TakeFrom/>
            )
            break;
        case "swap-discard":
            currentGameState = (
                <SwapDiscard/>
            )
            break;
        case "force-swap":
            currentGameState = (
                <ForceSwap/>
            )
            break;
        case "flip":
            currentGameState = (
                <div>Flip</div>
            )
        break;
        default:
            currentGameState = (
                <div>Bad Game State...</div>
            )
    }

    return (
        <div className="GameWindow">
            <Telemetry 
                playerServerId={playerId}
                gameState={player.turnType}
                cardInHandValue={data.cardInHand?.value ?? null}
                lastDiscardedCardValue={data.lastDiscardedCard?.value ?? null}
            />
            <div className="GameWindow_DeckWrapper">
                <Pile headerLabel="Discard Pile" faceUpCard={data.lastDiscardedCard} hasDeck={true}/>
                <CardGrid hand={player.hand}/>
                <Pile headerLabel="Card in Hand" faceUpCard={data.cardInHand} hasDeck={false}/>
                <ScoreLabel/>
            </div>
            <div className="GameWindow_ButtonWrapper">
                {currentGameState}
            </div>
        </div>
    )
}

export default GameWindow