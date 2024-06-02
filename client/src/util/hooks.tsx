import { useData } from "../ui/LobbyScreen";
import { IPlayer } from "./interfaces";

// Access to player - also checks for if its their turn
export function usePlayer(): IPlayer {
    const data = useData();

    return {
        player: data?.players[data.playerId],
        isCurrentTurn: data?.playerId === data?.currentPlayerId
    };
}