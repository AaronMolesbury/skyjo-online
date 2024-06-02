import { useEffect, useState } from "react";
import { usePlayer } from "../util/hooks";
import { useData } from "./LobbyScreen";
import "../css/GameWindow.css";
import CardGrid from "./CardGrid";
import ScoreLabel from "./ScoreLabel";
import Piles from "./Piles";
import CardInHand from "./CardInHand";
import EndGameScreen from "./EndGameScreen";

function GameWindow() {
    const data = useData();
    const {player} = usePlayer();
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

    if (!player) {
        return;
    }

    if (data?.winnerId !== -1) {
        return <EndGameScreen/>;
    }
    
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
        );
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
        );
    }

    return (
        <div className="GameWindow">
            {deckWrapper}
        </div>
    );
}

export default GameWindow