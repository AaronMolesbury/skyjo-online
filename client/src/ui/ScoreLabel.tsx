import { useData } from "./LobbyScreen"
import "../css/ScoreLabel.css"

function ScoreLabel() {
    const data = useData();
    const player = data?.players[data.playerId];

    return (
        <div className="ScoreLabel">
            <div className="ScoreLabel_Text">Score: {player?.score}</div>
        </div>
    );
}

export default ScoreLabel