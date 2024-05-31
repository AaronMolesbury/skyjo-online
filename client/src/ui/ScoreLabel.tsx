import { useData } from "../App"
import "../css/ScoreLabel.css"

function ScoreLabel() {
    const data = useData();

    if (!data) {
        return (
            <div>Bad Data @ ScoreLabel</div>
        )
    }

    const player = data.players[data.playerId]

    return (
        <div className="ScoreLabel">
            <div className="ScoreLabel_Text">Score: {player.score}</div>
        </div>
    )
}

export default ScoreLabel