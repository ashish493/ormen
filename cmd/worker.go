package cmd

import (
	"fmt"
	"log"

	"github.com/ashish493/ormen/sailor"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sailorCmd)
	sailorCmd.Flags().StringP("host", "H", "0.0.0.0", "Hostname or IP address")
	sailorCmd.Flags().IntP("port", "p", 5556, "Port on which to listen")
	sailorCmd.Flags().StringP("name", "n", fmt.Sprintf("worker-%s", uuid.New().String()), "Name of the worker")
	sailorCmd.Flags().StringP("dbtype", "d", "memory", "Type of datastore to use for tasks (\"memory\" or \"persistent\")")
}

var sailorCmd = &cobra.Command{
	Use:   "sailor",
	Short: "Worker command to operate a Cube worker node.",
	Long: `cube worker command.

The worker runs tasks and responds to the manager's requests about task state.`,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		name, _ := cmd.Flags().GetString("name")
		dbType, _ := cmd.Flags().GetString("dbtype")

		log.Println("Starting worker.")
		w := sailor.New(name, dbType)
		api := sailor.Api{Address: host, Port: port, Worker: w}
		go w.RunTasks()
		go w.CollectStats()
		go w.UpdateTasks()
		log.Printf("Starting worker API on http://%s:%d", host, port)
		api.Start()
	},
}
