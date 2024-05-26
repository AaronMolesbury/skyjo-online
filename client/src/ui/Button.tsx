import { IButtonProps } from "../util/interfaces"
import "../css/Button.css"

function Button(props: IButtonProps) {
    return (
        <div className="Button" onClick={props.clickHandler}>
            <div className="Button_Label">{props.labelText}</div>
        </div>
    )
}

export default Button