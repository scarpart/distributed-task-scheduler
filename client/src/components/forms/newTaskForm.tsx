import React, { useState } from 'react'; 

const NewTaskForm: React.FC = () => {
	const [task, setTask] = useState({
		Status: "",
		CpuUsage: 0,
		MemUsage: 0,
		TaskName: "",
		TaskDescription: "",
		TaskID: 0,
		UserID: 0,
		ScheduledTime: "",
		// ...
	});

	const handleSubmit = (event: React.FormEvent) => {
		event.preventDefault();
		// TODO: make POST request to the server here.
	};

	const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		setTask({ ...task, [event.target.name]: event.target.value });
	};

	return (
    <form onSubmit={handleSubmit}>
      <label>
        Task Name:
        <input
          type="text"
          name="TaskName"
          value={task.TaskName}
          onChange={handleChange}
        />
      </label>
      <label>
        Task Description:
        <input
          type="text"
          name="TaskDescription"
          value={task.TaskDescription}
          onChange={handleChange}
        />
      </label>
		{/* TODO: insert other fields here. just a demonstration for now. */ }
      <input type="submit" value="Submit" />
    </form>
  );
}; 

export default NewTaskForm;


