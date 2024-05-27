import { useData } from "../App"
import "../css/ScoreLabel.css"

function ScoreLabel() {
    const data = useData();

    if (!data) {
        return (
            <div>Bad Data @ ScoreLabel</div>
        )
    }

    return (
        <div className="ScoreLabel">
            <div className="ScoreLabel_Text">Score: {data.score}</div>
        </div>
    )
}

export default ScoreLabel