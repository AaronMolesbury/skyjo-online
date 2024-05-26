import { Hand, Card } from "./types"

export interface ISocketData {
    lastDiscardedCard: Card,
    cardInHand: Card,
    gameState: string,
    playerId: number,
    players: {
        hand: Hand
    }[],
    cardClickEnabled: boolean
}

export interface ICardGridProps {
    hand: Hand
}

export interface IButtonProps {
    clickHandler: React.MouseEventHandler<HTMLDivElement> | undefined,
    labelText: string
}

export interface ITelemetryProps {
    lastDiscardedCardValue: number | null,
    cardInHandValue: number | null,
    gameState: string,
    playerServerId: number
}