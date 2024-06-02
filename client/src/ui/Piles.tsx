import { usePlayer } from "../util/hooks";
import { useWebSocket } from "./LobbyScreen";
import { IFaceUpCardProps } from "../util/interfaces";
import { cardColorLookup } from "../util/cardColors";
import "../css/Piles.css"
import Header from "./Header";
import DiscardButton from "./DiscardButton";

function Piles(props: IFaceUpCardProps) {
    const ws = useWebSocket();
    const {player, isCurrentTurn} = usePlayer();

    const disabled = (!isCurrentTurn || player?.turnType !== "take-from");

    const deckClicked = () => {
        if (!disabled) {
            ws?.send("take-from-deck");
        }
    }

    const deck = (
        <div className={disabled ? "Piles_Deck disabled" : "Piles_Deck"} onClick={deckClicked}>
                <div className="Piles_DeckCard"></div>
                <div className="Piles_DeckCard"></div>
                <div className="Piles_DeckCard"></div>
                <div className="Piles_DeckCard"></div>
                <div className="Piles_DeckCard"></div>
        </div>
    );

    const discardPileClicked = () => {
        if (!disabled) {
            ws?.send("take-from-discard");
        }
    }

    const faceUpCard = props.faceUpCard ? (
        <div 
            className={disabled ? "Piles_FaceUpCard disabled" : "Piles_FaceUpCard"} 
            onClick={discardPileClicked}
            style={(props.faceUpCard && props.faceUpCard.value !== null) ? {background: cardColorLookup[props.faceUpCard.value]} : {}}
        >
            {props.faceUpCard.value}
        </div>
    ) : null;

    const discardPile = (
        <div className={disabled ? "Piles_Deck disabled" : "Piles_Deck"}>
                    <div className="Piles_DeckCard"></div>
                    <div className="Piles_DeckCard"></div>
                    <div className="Piles_DeckCard"></div>
                    <div className="Piles_DeckCard"></div>
                    {faceUpCard ? faceUpCard : <div className="Piles_DeckCard"></div>}
        </div>
    );

    return (
        <div className="Piles">
            <Header text={"Deck"}/>
            {deck}
            <Header text={"Discard Pile"}/>
            {discardPile}
            <DiscardButton/>
        </div>
    );
}

export default Piles