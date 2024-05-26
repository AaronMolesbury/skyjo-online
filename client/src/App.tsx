import { createContext, useContext, useEffect, useState } from 'react'
import './App.css'
import { ISocketData } from './util/interfaces';
import GameWindow from './ui/GameWindow';

const WS_URL = "ws://localhost:8080";

const WebSocketContext = createContext<WebSocket|null>(null);
const DataContext = createContext<ISocketData|null>(null);

export function useWebSocket() {
  return useContext(WebSocketContext);
}

export function useData() {
  return useContext(DataContext);
}

function App() {
  const [data, setData] = useState<ISocketData|null>(null);
  const [socket, setSocket] = useState<WebSocket|null>(null)

  useEffect(() => {
    // Create WebSocket connection.
    const ws = new WebSocket(WS_URL);
    setSocket(ws);

    // Listen for messages - update data in GameState context
    ws.addEventListener("message", (event) => {
      setData(JSON.parse(event.data as string));
    });

    return () => {
      ws.close()
    }
  }, [])

  return (
    <WebSocketContext.Provider value={socket}>
      <DataContext.Provider value={data}>
        <GameWindow/>
      </DataContext.Provider>
    </WebSocketContext.Provider>
  )
}

export default App
