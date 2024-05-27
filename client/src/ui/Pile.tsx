import "../css/Pile.css"
import Header from "./Header"
import { IPileProps } from "../util/interfaces";
import { cardColorLookup } from "../util/cardColors";

function Pile(props: IPileProps) {

    const deck = (
        <div className="Pile_Deck">
                <div className="Pile_DeckCard"></div>
                <div className="Pile_DeckCard"></div>
                <div className="Pile_DeckCard"></div>
                <div className="Pile_DeckCard"></div>
                <div className="Pile_DeckCard"></div>
        </div>
    )
    return (
        <div className="Pile">
            <Header text={props.headerLabel}></Header>
            {props.hasDeck ? deck : <></>}
            <div 
            className="Pile_FaceUpCard" 
            style={
                (props.faceUpCard && props.faceUpCard.value !== null) ? {background: cardColorLookup[props.faceUpCard.value]} : {}
            }>
                {props.faceUpCard.value}
            </div>
        </div>
    )
}

export default Pile