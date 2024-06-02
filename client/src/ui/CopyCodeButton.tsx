import { useState } from "react";
import { useData } from "./LobbyScreen";
import "../css/CopyCodeButton.css";

function CopyCodeButton() {
    const data = useData();
    const [copied, setCopied] = useState<boolean>(false);

    const copyClicked = () => {
        navigator.clipboard.writeText(data?.lobbyCode + "");
        setCopied(true);
    }

    return (
        <>
            {copied ? <div className="CopyCodeButton_Text">Link Copied!</div> : <div className="CopyCodeButton" onClick={copyClicked}></div>}
        </>
    );
}

export default CopyCodeButton