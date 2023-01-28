import {FaTrashAlt} from "react-icons/fa"
import {BiCheckbox, BiCheckboxSquare} from "react-icons/bi"

import "./component_styling.css"

const Task = ({id, task, task_status, update, undo, deleteTask}: {id: string, task: string, task_status: boolean, update : () => void, undo : () => void, deleteTask : () => void}) => {    
    return (
        <div className='TaskRow'>
            {task_status ? 
            <button onClick={undo}>
                <BiCheckboxSquare />
            </button>
            :
            <button  onClick={update}>
                <BiCheckbox />
            </button>
            }
            <h2 className={`Task${task_status ? " strike" : ""}`}>{task}</h2>


            <button className="trash" onClick={deleteTask}>
                <FaTrashAlt />
            </button>
        </div>
    )

}

export default Task


