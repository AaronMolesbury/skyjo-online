export interface ISocketData {
    playerId: number,
    players: {
        hand: ({value: number | null}|null)[][]
    }[]
}