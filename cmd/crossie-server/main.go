package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gotwarlost/crossies/internal/server"
	"github.com/spf13/cobra"
)

func main() {
	var port int
	var root string
	cmd := &cobra.Command{
		Use:   "crossie-server",
		Short: "run a fully contained crossie server for development use",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			mux, err := server.Handler(root)
			if err != nil {
				return err
			}
			log.Println("start server on port:", port, ", using root:", root)
			return http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), mux)
		},
	}
	f := cmd.Flags()
	f.IntVarP(&port, "port", "p", 8989, "port to run server on")
	f.StringVar(&root, "root", server.DefaultRoot(), "root directory for static files")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
