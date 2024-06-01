import { createContext, useContext, useEffect, useState } from 'react'
import { ISiteConfig, ISocketData, ILobby } from '../util/interfaces'
import SITE_CONFIG from './SiteConfig'
import ReadyUpButton from './ReadyUpButton';
import GameWindow from './GameWindow';

const SiteConfigContext = createContext<ISiteConfig>(SITE_CONFIG)
const WebSocketContext = createContext<WebSocket|null>(null);
const DataContext = createContext<ISocketData|null>(null);

export function useSiteConfig() {
  return useContext(SiteConfigContext);
}

export function useWebSocket() {
  return useContext(WebSocketContext);
}

export function useData() {
  return useContext(DataContext);
}

function LobbyScreen(props: ILobby) {
    const siteConfig = useSiteConfig()

    const [data, setData] = useState<ISocketData|null>(null);
    const [socket, setSocket] = useState<WebSocket|null>(null)

    useEffect(() => {
        // Create WebSocket connection.
        const socketUrl = props.lobbyCode ? siteConfig.socketUrl + "/" + props.lobbyCode : siteConfig.socketUrl;
        const ws = new WebSocket(socketUrl);
        setSocket(ws);

        // Listen for messages - update data in GameState context
        ws.addEventListener("message", (event) => {
            const rawData = JSON.parse(event.data as string)
            setData(rawData);
            if (window.location.pathname === "/") {
                history.pushState(null, "", "/" + rawData.lobbyCode);
            }
        });

        return () => {
            ws.close()
        }
    }, [])

    return (
        <WebSocketContext.Provider value={socket}>
            <DataContext.Provider value={data}>
                {data?.started ? <GameWindow/> : (
                    <>
                        <div>Lobby: {data?.lobbyCode}</div>
                        <ReadyUpButton/>
                    </>
                )}
            </DataContext.Provider>
        </WebSocketContext.Provider>
    )
}

export default LobbyScreen