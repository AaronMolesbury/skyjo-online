import { useState } from 'react'
import './App.css'
import JoinLobbyForm from './ui/JoinCreateLobbyForm'
import LobbyScreen from './ui/LobbyScreen'

const pathRegex = /^\/\d{6}$/

function App() {
  const [lobbyCode, setLobbyCode] = useState<number | null>(null)
  const [gameJoined, setGameJoined] = useState(false)
  const [showTitle, setShowTitle] = useState(true)

  if (!gameJoined && window.location.pathname.match(pathRegex)) {
    const lc = Number(window.location.pathname.substring(1))
    setLobbyCode(lc)
    setGameJoined(true)
  } else {
    history.pushState(null, "", "/");
  }

  return (
    <div className="App">
        {showTitle ? <div className="App_Title">SKYJO</div> : <></>}
        {gameJoined ? <LobbyScreen lobbyCode={lobbyCode} setShowTitleCallback={setShowTitle}/> : <JoinLobbyForm setGameJoinedCallback={setGameJoined} setLobbyCodeCallback={setLobbyCode}/>}
    </div>
  )
}

export default App
