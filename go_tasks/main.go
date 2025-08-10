package main

import (
	"encoding/json"
	"fmt"
	"os"
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
				// Open JSON file
				fileHandle, err = os.OpenFile("tasks.json", os.O_APPEND, 0644)
				if err != nil {
					fmt.Println("Error opening file: ", err)
					return
				}

				// Decode JSON file
				fileContents, err := os.ReadFile("tasks.json")
				if err != nil {
					fmt.Println("Error opening tasks JSON: ", err)
					return
				}

				if len(fileContents) != 0 {
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
