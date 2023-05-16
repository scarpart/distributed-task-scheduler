import React, { useEffect, useState } from 'react';

type Task = {
  Status: string;
  CpuUsage: number;
  MemUsage: number;
  TaskName: string;
  TaskDescription: string;
  TaskID: number;
  UserID: number;
  ScheduledTime: string;
  NThreads: number;
  Priority: number;
  RetryCount: number;
  MaxRetries: number;
  Dependencies: number[];
  CreatedTime: string;
  LastUpdatedTime: string;
  ErrorMessage: string;
  MachineID: string;
  Output: string;
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
		  <tbody>
			{data.map((task) => (
			  <tr key={task.TaskID}>
			  <th>{task.TaskID}</th>
			  <th>{task.MachineID}</th>
			  <th>{task.TaskDescription}</th>
			  <th>{task.Status}</th>
			  <th>{task.UserID}</th>
			  <th>{task.ScheduledTime}</th>
			  <th>{task.LastUpdatedTime}</th>
			  <th>{task.Priority}</th>
			  <th>{task.Dependencies}</th>
			  </tr>
			))}
		  </tbody>
		</thead>
	  </table>
    </div>
  );
};

export default App;
