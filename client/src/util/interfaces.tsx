import { Hand, Card } from "./types"

export interface ISiteConfig {
    socketUrl: string
}

export interface ILobby {
    lobbyCode: number | null
}

export interface IJoinCreateLobbyForm {
    setLobbyCodeCallback: React.Dispatch<React.SetStateAction<number|null>>
    setGameJoinedCallback: React.Dispatch<React.SetStateAction<boolean>>
}

export interface ISocketData {
    started: boolean,
    lobbyCode: number,
    lastDiscardedCard: Card | null,
    cardInHand: Card | null,
    playerId: number,
    players: ISocketPlayer[],
    currentPlayerId: number,
    winnerId: number
}

export interface ISocketPlayer {
    hand: Hand,
    turnType: string,
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
    faceUpCard: Card | null,
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
    labelText: string,
    className?: string
}
