import React, {Component} from 'react'
import axios from "axios";
import Task from "../Components/Task"

import './list_styling.css'

let endpoint = "http://10.0.0.6:9000"




class Home extends Component<any, any> {


  constructor(props: any){
    super(props)

    this.state = {
      task:"",
      items:[]
    }
  }
  
  getTask = () => {
    axios.get(endpoint + "/api/tasks").then((res)=>{
      if (res.data) {
        this.setState({
          items: res.data.map((task: any)=>{
            return (

              <Task 
              id={task._id}
              task={task.task}
              task_status={task.status}
              update={() => {this.updateTask(task._id)}}
              undo={() => {this.undoTask(task._id)}}
              deleteTask={() => {this.deleteTask(task._id)}}
              key={task._id}
              />
            );
          }),
        });
      }
      else {
        this.setState({
          items: []
        });
      }
    });
  };

  updateTask = (id : string) => {
    axios.put(endpoint + "/api/tasks/" + id).then(() => {
      this.getTask()
    })

  }
  undoTask = (id : string) => {
    axios.put(endpoint + "/api/undoTask/" + id).then(() => {
      this.getTask()
    })
  }

  deleteTask = (id : string) => {
    axios.delete(endpoint + "/api/deleteTask/" + id).then(() => {
      this.getTask()
    })
  }

  deleteAll = () => {
    axios.delete(endpoint + "/api/deleteAllTasks").then(() => {this.getTask()})
  }

  onChange = (e: any) => {
    this.setState({
      [e.target.name] : e.target.value
    })
  }

  onSubmit = () => {
    let {task} = this.state;

    console.log(task)
    
    if (task) {
      axios.post(endpoint + "/api/task", 
      {task: task,})
      .then((res)=> {
        this.getTask();
        this.setState({
          task:"",
        });
        console.log(res)
      });
    }
  };


  componentDidMount() {
      this.getTask()
  }
  render() {
    return (
      <div>
          <h2>
            To Do List
          </h2>
        <div className='row'>
          <div>
            <input
            type="text"
            name="task"
            className='taskInput'
            onChange={this.onChange}
            value={this.state.task}
            placeholder="Create A Task"
            />
          </div>

          <button className='TaskButton' onClick={this.onSubmit}>Submit</button>

          <button className='TaskButton' onClick={this.deleteAll}>Delete</button>
        </div>
        <div className='Tasks'>
          {this.state.items}
        </div>
      </div>
    )
  }
}

export default Home