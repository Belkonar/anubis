package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/Belkonar/anubis/types"
	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
)

var routers = map[string]*chi.Mux{}
var cfgFile string

func globalHandler(w http.ResponseWriter, r *http.Request) {
	uriParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")

	prefix := uriParts[0]
	uriParts = uriParts[1:]
	newPath := "/" + strings.Join(uriParts, "/")

	fmt.Println(prefix, newPath)

	r.URL.Path = newPath // Reset path to remove prefix

	router, ok := routers[prefix]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No router found for %s", prefix)
		return
	}

	router.ServeHTTP(w, r)
}

func makeRouters() {
	fmt.Println(cfgFile)

	if cfgFile == "" {
		panic("No config file specified")
	}

	configData, err := os.ReadFile(cfgFile)

	if err != nil {
		panic(err)
	}

	config := []types.TargetConfig{}

	json.Unmarshal(configData, &config)

	for _, target := range config {
		// Fallback to proxy
		setupRouter(target)
	}
}

func setupRouter(target types.TargetConfig) {
	router := chi.NewRouter()

	proxy := makeProxy(target.Target)

	router.Get("/{asd}", func(w http.ResponseWriter, r *http.Request) {
		asd := chi.URLParam(r, "asd")
		//fmt.Println(chi.RouteContext(r.Context()).URLParams.Keys)
		if asd != "asd" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Lol nope")
			return
		}

		proxy.ServeHTTP(w, r)
	})

	router.NotFound(proxy.ServeHTTP) // Catch all router

	routers[target.Prefix] = router
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
		makeRouters()
		globalHandler := http.HandlerFunc(globalHandler)

		err := http.ListenAndServe("127.0.0.1:8000", globalHandler)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gatewayCmd)

	gatewayCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "Config File")
}
