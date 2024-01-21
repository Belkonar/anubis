/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/cobra"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Echoing %s", r.URL.Path)
}

// echoerCmd represents the echoer command
var echoerCmd = &cobra.Command{
	Use:   "echoer",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		router := chi.NewRouter()
		router.Use(middleware.Logger)
		router.NotFound(echoHandler)

		err := http.ListenAndServe("127.0.0.1:8010", router)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(echoerCmd)
}
