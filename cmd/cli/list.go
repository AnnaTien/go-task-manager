package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go-task-manager/internal/task"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get("http://localhost:8080/tasks")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var tasks []task.Task
			if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
				log.Fatal(err)
			}
			for _, t := range tasks {
				fmt.Printf("ID: %d, Name: %s, Completed: %t\n", t.ID, t.Name, t.Completed)
			}
		} else {
			fmt.Printf("Failed to get tasks. Status: %s\n", resp.Status)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
