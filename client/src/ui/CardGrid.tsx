import { ICardGridProps } from "../util/interfaces";
import "../css/CardGrid.css";
import Card from "./Card.tsx";

function CardGrid(props: ICardGridProps) {
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
    );
}

export default CardGrid