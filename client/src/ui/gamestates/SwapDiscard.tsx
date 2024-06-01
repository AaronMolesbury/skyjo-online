import { useWebSocket } from "../LobbyScreen"
import Button from "../Button"

function SwapDiscard() {
    const ws = useWebSocket();
    if (!ws) {
        return (
            <div>Disconnected...</div>
        )
    }

    const discardClicked = () => {
        ws.send("discard")
    }

    return (
        <>
            <div>Swap or </div>
            <Button clickHandler={discardClicked} labelText="Discard & Flip"/>
        </>
    )
}

export default SwapDiscard