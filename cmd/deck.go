package cmd

import (
	"log"

	"github.com/ashish493/ormen/deck"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deckCmd)
	deckCmd.Flags().StringP("host", "H", "0.0.0.0", "Hostname or IP address")
	deckCmd.Flags().IntP("port", "p", 5555, "Port on which to listen")
	deckCmd.Flags().StringSliceP("workers", "w", []string{"localhost:5556"}, "List of workers on which the manger will schedule tasks.")
	deckCmd.Flags().StringP("scheduler", "s", "epvm", "Name of scheduler to use.")
	deckCmd.Flags().StringP("dbType", "d", "memory", "Type of datastore to use for events and tasks (\"memory\" or \"persistent\")")
}

var deckCmd = &cobra.Command{
	Use:   "deck",
	Short: "command to operate a manager node.",
	Long: `cube manager command.

The manager controls the orchestration system and is responsible for:
- Accepting tasks from users
- Scheduling tasks onto worker nodes
- Rescheduling tasks in the event of a node failure
- Periodically polling workers to get task updates`,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		workers, _ := cmd.Flags().GetStringSlice("workers")
		scheduler, _ := cmd.Flags().GetString("scheduler")
		dbType, _ := cmd.Flags().GetString("dbType")

		log.Println("Starting manager.")
		m := deck.New(workers, scheduler, dbType)
		api := deck.Api{Address: host, Port: port, Manager: m}
		go m.ProcessTasks()
		go m.UpdateTasks()
		go m.DoHealthChecks()
		go m.UpdateNodeStats()
		log.Printf("Starting manager API on http://%s:%d", host, port)
		api.Start()
	},
}
