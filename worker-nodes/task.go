package main

type Task struct {
	ID     		string  			`json:"id"`
	Image  		string				`json:"image"`
	Command 	string 				`json:"command"`
	Args    	[]string			`json:"args"`

	// Environment is a map of the env variables and their values that should be used inside the container 
	Environment map[string]string   `json:"environment"`
}



