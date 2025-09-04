package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go-task-manager/internal/task"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [task name]",
	Short: "Add a new task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskName := args[0]
		newTask := task.Task{Name: taskName}

		jsonBody, err := json.Marshal(newTask)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.Post("http://localhost:8080/tasks", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			fmt.Println("Task added successfully.")
		} else {
			fmt.Printf("Failed to add task. Status: %s\n", resp.Status)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
