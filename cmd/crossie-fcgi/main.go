package main

import (
	"net/http/fcgi"
	"os"

	"github.com/gotwarlost/crossies/internal/server"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "api.fcgi",
		Short: "run a FastCGI version of the crossie API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			mux, err := server.CGIHandler()
			if err != nil {
				return err
			}
			return fcgi.Serve(nil, mux)
		},
	}
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
