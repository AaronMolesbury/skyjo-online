import { BaseSyntheticEvent } from "react";
import { IJoinCreateLobbyForm } from "../util/interfaces";
import "../css/JoinCreateLobbyForm.css";
import Button from "./Button";

function JoinCreateLobbyForm(props: IJoinCreateLobbyForm) {
    const validateLength = (e: BaseSyntheticEvent) => {
        if (e.target.value.length > e.target.maxLength) {
            e.target.value = e.target.value.substring(0, 6);
        } 
    }

    const joinClicked = (e: BaseSyntheticEvent) => {
        const inputTarget = e.target[0];
        if (inputTarget.value.length === inputTarget.maxLength) {
            history.pushState(null, "", "/" + inputTarget.value);
        }
        props.setLobbyCodeCallback(inputTarget.value);
        props.setGameJoinedCallback(true);
    }

    const createClicked = () => {
        props.setGameJoinedCallback(true);
    }

    return (
        <>
            <form onSubmit={joinClicked} className="JoinCreateLobbyForm">
                    <input 
                        className="JoinCreateLobbyForm_Input" 
                        type="number" 
                        required minLength={6} 
                        maxLength={6} 
                        placeholder="Enter a 6 digit lobby code" 
                        onInput={validateLength}
                    />
                    <input
                        className="JoinCreateLobbyForm_Submit" type="submit" value="Join"
                    />
            </form>
            <div className="JoinCreateLobbyForm_Text">Or</div>
            <Button className="JoinCreateLobbyForm_Button" clickHandler={createClicked} labelText={"Create Lobby"}/>
        </>
    );
}

export default JoinCreateLobbyForm