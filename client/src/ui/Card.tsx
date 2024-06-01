import "../css/Card.css"
import { useWebSocket, useData } from "./LobbyScreen"
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
    const player = data.players[data.playerId]

    let className: string;
    if ((value !== null && FACE_UP_HOVER_ENABLED_GAMESTATES.includes(player.turnType)) || (value === null && FACE_DOWN_HOVER_ENABLED_GAMESTATES.includes(player.turnType))) {
        className = "Card";
    } else {
        className = "Card disabled";
    }

    const style = value !== null ? {background: cardColorLookup[value]} : {};

    return (
        <div 
        className={className}
        key={props.colIndex} 
        style={style} 
        onClick={() => {
            //Stop flips on face up cards
            if ((player.turnType === "flip") && value) {
                return;
            }

            ws.send(`card-clicked;${props.colIndex};${props.rowIndex}`)
        }}>
            {value}
        </div>
    )
}

export default Card