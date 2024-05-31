import "../css/Pile.css"
import Header from "./Header"
import { IPileProps } from "../util/interfaces";
import { cardColorLookup } from "../util/cardColors";

function Pile(props: IPileProps) {
    if (!props.faceUpCard && !props.hasDeck) {
        // Reserves space for component
        return <div className="Pile"></div>
    }

    let deck = <></>
    if (props.hasDeck) {
        deck = (
            <div className="Pile_Deck">
                    <div className="Pile_DeckCard"></div>
                    <div className="Pile_DeckCard"></div>
                    <div className="Pile_DeckCard"></div>
                    <div className="Pile_DeckCard"></div>
                    <div className="Pile_DeckCard"></div>
            </div>
        )
    }

    let card = <></>
    if (props.faceUpCard) {
        card = (
            <div 
            className="Pile_FaceUpCard" 
            style={
                (props.faceUpCard && props.faceUpCard.value !== null) ? {background: cardColorLookup[props.faceUpCard.value]} : {}
            }>
                {props.faceUpCard.value}
            </div>
        )
    }

    return (
        <div className="Pile">
            <Header text={props.headerLabel}></Header>
            {deck}
            {card}
        </div>
    )
}

export default Pile