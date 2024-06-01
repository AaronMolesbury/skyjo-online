import { ICardGridProps } from "../util/interfaces"
import "../css/CardGrid.css"
import { useWebSocket, useData } from "./LobbyScreen"
import Card from "./Card.tsx"

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
                        <Card card={card} colIndex={colIndex} rowIndex={rowIndex}/>
                    ))}
                </div>
            ))}
        </div>
    )
}

export default CardGrid