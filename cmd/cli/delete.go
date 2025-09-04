package cli

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [task ID]",
	Short: "Delete a task by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/tasks/%s", taskID), nil)
		if err != nil {
			log.Fatal(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent {
			fmt.Printf("Task with ID %s deleted successfully.\n", taskID)
		} else {
			fmt.Printf("Failed to delete task. Status: %s\n", resp.Status)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
