import { createContext, useContext, useEffect, useState } from 'react'
import { ISiteConfig, ISocketData, ILobby } from '../util/interfaces'
import SITE_CONFIG from './SiteConfig'
import ReadyUpButton from './ReadyUpButton';
import GameWindow from './GameWindow';
import "../css/LobbyScreen.css"

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
    const [copied, setCopied] = useState<boolean>(false);

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

    const copyClicked = () => {
        if (socket && data) {
            navigator.clipboard.writeText(data?.lobbyCode + "");
            setCopied(true)
        }
    }

    let content: JSX.Element;
    if (!socket || !data)  {
        content = <div>Failed to create lobby - Please refresh</div>
    } else {
        if (data.started) {
            props.setShowTitleCallback(false)
            content = <GameWindow/>
        } else {
            content =  (
                <>
                    <div className="LobbyScreen_InfoWrapper">
                        <div className="LobbyScreen_LobbyText">Lobby: {data?.lobbyCode}</div>
                        {copied ? <div className="LobbyScreen_CopiedText">Link Copied!</div> : <div className="LobbyScreen_CopyURL" onClick={copyClicked}></div>}
                    </div>
                    <ReadyUpButton/>
                </>
            )
        }
    }

    return (
        <WebSocketContext.Provider value={socket}>
            <DataContext.Provider value={data}>
                {content}
            </DataContext.Provider>
        </WebSocketContext.Provider>
    )
}

export default LobbyScreen