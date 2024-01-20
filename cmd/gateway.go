/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
)

func globalHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	router := chi.NewRouter()

	proxy := makeProxy("http://example.com")

	router.Get("/", proxy.ServeHTTP)

	// Fallback to proxy
	router.NotFound(proxy.ServeHTTP)

	router.ServeHTTP(w, r)

	// proxy.ServeHTTP(w, r)
}

func makeProxy(target string) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{ // TODO: Make a factory so this can be reused
		Rewrite: makeRewriter(target),
	}
}

func makeRewriter(target string) func(*httputil.ProxyRequest) {
	return func(r *httputil.ProxyRequest) {
		target, err := url.Parse(target)

		if err != nil {
			panic(err) // FIXME: handle error
		}

		r.SetURL(target)

		r.SetXForwarded()
		r.Out.Host = target.Host // Super annoying but entirely necessary
		fmt.Println(r.Out.URL)
	}
}

// gatewayCmd represents the gateway command
var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		globalHandler := http.HandlerFunc(globalHandler)

		err := http.ListenAndServe("127.0.0.1:8000", globalHandler)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gatewayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gatewayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gatewayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
