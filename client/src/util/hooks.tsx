import { useData } from "../ui/LobbyScreen";
import { IPlayer } from "./interfaces";

export function usePlayer(): IPlayer {
    const data = useData();

    return {
        player: data?.players[data.playerId],
        isCurrentTurn: data?.playerId === data?.currentPlayerId
    }
}