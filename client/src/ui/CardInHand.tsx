import { IFaceUpCardProps } from "../util/interfaces";
import { cardColorLookup } from "../util/cardColors";
import Header from "./Header";
import "../css/CardInHand.css";

function CardInHand (props: IFaceUpCardProps) {

    const card = props.faceUpCard ? (
        <div 
            className="CardInHand_Card"
            style={(props.faceUpCard && props.faceUpCard.value !== null) ? {background: cardColorLookup[props.faceUpCard.value]} : {}}
        >
            {props.faceUpCard.value}
        </div>
    ) : <div className="CardInHand_DefaultCard">?</div>;

    return (
        <div className="CardInHand">
            <Header text={"Card In Hand"}/>
            {card}
        </div>
    )
}

export default CardInHand