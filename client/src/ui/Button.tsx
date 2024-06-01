import { IButtonProps } from "../util/interfaces"
import "../css/Button.css"

function Button(props: IButtonProps) {
    let className = "Button ";
    if (props.className) {
        className += props.className;
    }

    return (
        <div className={className} onClick={props.clickHandler}>
            <div className="Button_Label">{props.labelText}</div>
        </div>
    )
}

export default Button