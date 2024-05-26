import { ITelemetryProps } from "../util/interfaces"
import "../css/Telemetry.css"

function Telemetry(props: ITelemetryProps) {
    return (
        <div className="Telemetry">
            <div className="Telemetry_Data">Server Player ID: {props.playerServerId}</div>
            <div className="Telemetry_Data">Game State: {props.gameState}</div>
            <div className="Telemetry_Data">Card in Hand: {props.cardInHandValue}</div>
            <div className="Telemetry_Data">Last Discarded Card: {props.lastDiscardedCardValue}</div>
        </div>
    )
}

export default Telemetry