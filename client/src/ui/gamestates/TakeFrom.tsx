import { useWebSocket } from "../LobbyScreen";
import Button from "../Button"

function TakeFrom() {
    const ws = useWebSocket();

    if (!ws) {
        return (
            <div>Disconnected...</div>
        )
    }

    const takeFromDeckClicked = () => {
        ws.send("take-from-deck");
    }

    const takeFromDiscardPileClicked = () => {
        ws.send("take-from-discard");
    }

    return (
        <>
            <Button clickHandler={takeFromDeckClicked} labelText="Take from Deck"/>
            <Button clickHandler={takeFromDiscardPileClicked} labelText="Take from Discard Pile"/>
        </>
    )
}

export default TakeFrom