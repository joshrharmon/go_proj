package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

var validCommands = map[string]struct{}{
	"add":              {},
	"update":           {},
	"delete":           {},
	"mark-in-progress": {},
	"mark-done":        {},
	"list":             {},
}

type Task struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func main() {
	// Capture cli args
	cmdArgs := os.Args[1:]

	// Check that command is supplied
	if len(cmdArgs) == 0 {
		fmt.Println(":: No valid command supplied. Exiting...")
	} else {
		taskArg := cmdArgs[0]

		// Lookup the command safely
		_, ok := validCommands[taskArg]

		if ok {
			var fileHandle *os.File
			var jsonDecoder *json.Decoder
			var loadedTasks []Task

			// Store command into var
			var command = string(taskArg)

			// Check existence of JSON file
			_, err := os.Stat("tasks.json")

			// Create JSON if necessary
			if os.IsNotExist(err) {
				fileHandle, err = os.Create("tasks.json")
				if err != nil {
					fmt.Println("Error creating file: ", err)
					return
				}
			} else {
				// Decode JSON file
				fileContents, err := os.ReadFile("tasks.json")
				if err != nil {
					fmt.Println("Error opening tasks JSON: ", err)
					return
				}

				if len(fileContents) != 0 {
					// Open JSON file
					fileHandle, err = os.Open("tasks.json")
					if err != nil {
						fmt.Println("Error opening file: ", err)
						return
					}
					jsonDecoder = json.NewDecoder(fileHandle)
					if err := jsonDecoder.Decode(&loadedTasks); err != nil {
						fmt.Println("Error decoding JSON: ", err)
						return
					}
				}

				err = fileHandle.Close()
				if err != nil {
					fmt.Println("There was an error with closing the file handle: ", err)
					return
				}
			}

			switch command {
			case "add":
				add(cmdArgs, loadedTasks)
			case "update", "delete":
				update(cmdArgs, loadedTasks)
			case "list":
				list(loadedTasks)
			}

		} else {
			fmt.Printf("This was an invalid command %s\n", taskArg)
			return
		}
	}
}

func add(cmdArgs []string, loadedTasks []Task) {
	if len(cmdArgs) < 2 {
		fmt.Println("Supply a name for the task: add \"Name of task\"")
		return
	} else {
		var err error
		// Construct new task
		latestId := findLatestId(loadedTasks)
		task := Task{
			ID:        latestId + 1,
			Name:      cmdArgs[1],
			Status:    "todo",
			CreatedAt: getCurrentTimeString(),
			UpdatedAt: getCurrentTimeString(),
		}
		loadedTasks = append(loadedTasks, task)
		fileHandle, err := os.Create("tasks.json")
		if err != nil {
			fmt.Println("Error opening file for writing: ", err)
			return
		}
		defer fileHandle.Close()

		jsonEncoder := json.NewEncoder(fileHandle)
		if err := jsonEncoder.Encode(loadedTasks); err != nil {
			fmt.Println("Error writing JSON: ", err)
			return
		}

		if loadedTasks[len(loadedTasks)-1].ID == latestId+1 {
			fmt.Printf("Task \"%s\" successfully added.\n", task.Name)
		} else {
			fmt.Printf("There was an error adding task \"%s\".\n", task.Name)
		}
	}
}

func update(cmdArgs []string, loadedTasks []Task) {
	if len(cmdArgs) < 0 {
		fmt.Println("Error, usage: update [ID#] \"New task name\"")
		return
	} else {
		// Check for integer value in update command
		cmdArgCmd := cmdArgs[0]
		cmdArgId := cmdArgs[1]
		var cmdArgNewName string
		if len(cmdArgs) > 2 {
			cmdArgNewName = cmdArgs[2]
		}
		var err error
		var id int
		var indexOfTask int
		var taskToUpdate Task

		if id, err = strconv.Atoi(cmdArgId); err != nil {
			fmt.Printf("Please enter an ID as such: update [ID#] \"New task name\" OR delete [ID#]")
			return
		} else if taskToUpdate, indexOfTask, err = findTask(loadedTasks, id); err != nil {
			fmt.Printf("Error: %s", err)
			return
		}

		switch cmdArgCmd {
		case "update":
			task := Task{
				ID:        id,
				Name:      cmdArgNewName,
				Status:    taskToUpdate.Status,
				CreatedAt: taskToUpdate.CreatedAt,
				UpdatedAt: getCurrentTimeString(),
			}
			loadedTasks[indexOfTask] = task
		case "delete":
			loadedTasks = removeFromSlice(loadedTasks, indexOfTask)
		}

		fileHandle, err := os.Create("tasks.json")
		if err != nil {
			fmt.Println("Error opening file for writing: ", err)
			return
		}
		defer fileHandle.Close()

		jsonEncoder := json.NewEncoder(fileHandle)
		if err := jsonEncoder.Encode(loadedTasks); err != nil {
			fmt.Println("Error writing JSON: ", err)
			return
		}

		if cmdArgCmd == "update" {
			if loadedTasks[indexOfTask].Name == cmdArgNewName {
				fmt.Printf("Task successfully updated.\n")
			} else {
				fmt.Printf("There was an error updating task\n")
				return
			}
		} else if cmdArgCmd == "delete" {
			if _, i, _ := findTask(loadedTasks, id); i == -1 {
				fmt.Printf("Task \"%s\" successfully deleted.\n", cmdArgId)
			} else {
				fmt.Printf("Task with ID \"%s\" was not deleted due to an error.\n", cmdArgId)
			}
		}

	}
}

/*
 * Basic list command to show current tasks
 */
func list(loadedTasks []Task) {
	for _, task := range loadedTasks {
		fmt.Printf("ID: %d / Status: %s / Task: %s | Created: %s, Updated: %s\n",
			task.ID,
			task.Status,
			task.Name,
			task.CreatedAt,
			task.UpdatedAt,
		)
	}
}

// Will return the last ID
func findLatestId(loadedTasks []Task) int {
	var latestId int
	if len(loadedTasks) == 0 {
		latestId = 0
	} else {
		latestId = loadedTasks[len(loadedTasks)-1].ID
	}
	return latestId
}

// Removes element from slice
func removeFromSlice(slice []Task, indexToRemove int) []Task {
	return append(slice[:indexToRemove], slice[indexToRemove+1:]...)
}

// Find task by ID
func findTask(loadedTasks []Task, id int) (Task Task, index int, error error) {
	for i, task := range loadedTasks {
		if task.ID == id {
			return task, i, nil
		}
	}
	return Task, -1, fmt.Errorf("Task with ID %d was not found.", id)
}

func getCurrentTimeString() string {
	currentTime := time.Now()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second())
}
