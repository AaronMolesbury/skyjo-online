import { useWebSocket } from "./LobbyScreen";
import { usePlayer } from "../util/hooks";
import { cardColorLookup } from "../util/cardColors";
import { ICardProps } from "../util/interfaces";
import "../css/Card.css";

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
    const {player} = usePlayer();
    if (!ws || !player || !props.card) {
        return;
    }

    const value = props.card?.value ?? null;

    let className: string;
    if (
        (value !== null && FACE_UP_HOVER_ENABLED_GAMESTATES.includes(player.turnType)) ||
        (value === null && FACE_DOWN_HOVER_ENABLED_GAMESTATES.includes(player.turnType))
    ) {
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

            ws.send(`card-clicked;${props.colIndex};${props.rowIndex}`);
        }}>
            {value}
        </div>
    );
}

export default Card