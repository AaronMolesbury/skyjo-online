import { ICardGridProps } from "../util/interfaces"
import "../css/CardGrid.css"
import { useWebSocket, useData } from "../App"

function CardGrid(props: ICardGridProps) {
    const ws = useWebSocket();
    const data = useData();

    if (!ws || !data) {
        return (
            <div>Bad Connection/Data @ CardGrid</div>
        )
    }

    return (
        <div className="CardGrid">
            {props.hand.map((row, rowIndex) => (
                <div className="CardGrid_Row" key={rowIndex}>
                    {row.map((card, colIndex) => (
                        <div className="CardGrid_Card" key={colIndex} onClick={() => {
                            if (!data.cardClickEnabled) {
                                return;
                            }

                            // Stop flips on face up cards
                            if ((data.gameState === "flip" || data.gameState === "flip-two") && card?.value) {
                                return;
                            }

                            ws.send(`card-clicked;${colIndex};${rowIndex}`)
                        }}>
                            {card?.value}
                        </div>
                    ))}
                </div>
            ))}
        </div>
    )
}

export default CardGrid