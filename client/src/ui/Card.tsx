import "../css/Card.css"
import { useWebSocket, useData } from "../App"
import { cardColorLookup } from "../util/cardColors";
import { ICardProps } from "../util/interfaces"

const FACE_UP_HOVER_ENABLED_GAMESTATES = [
    "swap-discard",
    "force-swap"
];

const FACE_DOWN_HOVER_ENABLED_GAMESTATES = [
    "swap-discard",
    "force-swap",
    "flip",
    "flip-two",
];  

function Card(props: ICardProps) {
    const ws = useWebSocket();
    const data = useData();

    if (!ws || !data) {
        return (
            <div>Bad Connection/Data @ Card</div>
        )
    }

    if (!props.card) {
        return <></>
    }

    const value = props.card?.value ?? null;

    let className: string;
    if ((value && FACE_UP_HOVER_ENABLED_GAMESTATES.includes(data.gameState)) || (!value && FACE_DOWN_HOVER_ENABLED_GAMESTATES.includes(data.gameState))) {
        className = "Card";
    } else {
        className = "Card disabled";
    }

    return (
        <div 
        className={className}
        key={props.colIndex} 
        style={
            value ? {background: cardColorLookup[value]} : {}
        } 
        onClick={() => {
            //Stop flips on face up cards
            if ((data.gameState === "flip" || data.gameState === "flip-two") && value) {
                return;
            }

            ws.send(`card-clicked;${props.colIndex};${props.rowIndex}`)
        }}>
            {value}
        </div>
    )
}

export default Card