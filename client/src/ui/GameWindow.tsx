import { useData, useWebSocket } from "../App"
import "../css/GameWindow.css"
import Button from "./Button";
import CardGrid from "./CardGrid";
import Pile from "./Pile";
import ScoreLabel from "./ScoreLabel";
import Telemetry from "./Telemetry";
import FlipTwo from "./gamestates/FlipTwo";
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

    const playerId = data.playerId;
    const hand = data.players[playerId].hand;
    
    let currentGameState: JSX.Element;
    switch (data.gameState) {
        case "flip-two":
            currentGameState = (
                <FlipTwo/>
            )
            break;
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
                gameState={data.gameState}
                cardInHandValue={data.cardInHand.value}
                lastDiscardedCardValue={data.lastDiscardedCard.value}
            />
            <div className="GameWindow_DeckWrapper">
                <Pile headerLabel="Discard Pile" faceUpCard={data.lastDiscardedCard} hasDeck={true}/>
                <CardGrid hand={hand}/>
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