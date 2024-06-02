import { IHeaderProps } from "../util/interfaces"
import "../css/Header.css"

function Header(props: IHeaderProps) {
    return (
        <div className="Header">
            {props.text}
        </div>
    );
}

export default Header