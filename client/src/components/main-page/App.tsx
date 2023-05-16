import React, { useEffect, useState } from 'react';

type Task = {
  status: string;
  cpuUsage: number;
  memUsage: number;
  taskName: string;
  taskDescription: string;
  taskID: number;
  userID: number;
  scheduledTime: string;
  nThreads: number;
  priority: number;
  retryCount: number;
  maxRetries: number;
  dependencies: number[];
  createdTime: string;
  lastUpdatedTime: string;
  errorMessage: string;
  machineID: string;
  output: string;
};

const App: React.FC = () => {
  const [data, setData] = useState<Task[] | null>(null);
  const [error, setError] = useState<string | null>(null);
  
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch('https://127.0.0.1:9090/scheduled-tasks', {
          method: 'GET',
          mode: 'cors',
          cache: 'no-cache',
          credentials: 'same-origin',
          headers: {
            'Content-Type': 'application/json'
          },
          redirect: 'follow',
          referrerPolicy: 'no-referrer',
        });
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data: Task[] = await response.json();
        setData(data);
		console.log(data);
      } catch (error) {
		if (error instanceof Error) {
			setError(error.message);
		} else {
			setError("An unknown error occurred!");
		}
      }
    };

    fetchData();
  }, []);

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (!data) {
    return <div>Loading...</div>;
  }

  return (
    <div>
	  <table>
	  	<thead>
		  <tr>
			<th>Task ID</th>
			<th>Machine ID</th>
			<th>Task Description</th>
			<th>Status</th>
			<th>Owner</th>
			<th>Start Time</th>
			<th>Last Update Time</th>
			<th>Priority</th>
			<th>Dependencies</th>
		  </tr>
		</thead>
		<tbody>
		{data.map((task, index) => (
					<tr key={`${task.taskID}-${index}`}>
					<td>{task.taskID}</td>
					<td>{task.machineID}</td>
					<td>{task.taskDescription}</td>
					<td>{task.status}</td>
					<td>{task.userID}</td>
					<td>{task.scheduledTime}</td>
					<td>{task.lastUpdatedTime}</td>
					<td>{task.priority}</td>
					<td>{task.dependencies ? task.dependencies.join(", ") : "null"}</td>
					</tr>
					))}
  		</tbody>
	  </table>
    </div>
  );
};

export default App;
