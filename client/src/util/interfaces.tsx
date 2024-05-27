import { Hand, Card } from "./types"

export interface ISocketData {
    lastDiscardedCard: Card,
    cardInHand: Card,
    gameState: string,
    playerId: number,
    players: {
        hand: Hand
    }[],
    score: number
}

export interface ITelemetryProps {
    lastDiscardedCardValue: number | null,
    cardInHandValue: number | null,
    gameState: string,
    playerServerId: number
}

export interface IPileProps {
    headerLabel: string,
    faceUpCard: Card,
    hasDeck: boolean
}

export interface IHeaderProps {
    text: string
}

export interface ICardGridProps {
    hand: Hand
}

export interface ICardProps {
    card: Card | null,
    colIndex: number,
    rowIndex: number
}

export interface IButtonProps {
    clickHandler: React.MouseEventHandler<HTMLDivElement> | undefined,
    labelText: string
}
