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
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/cobra"
)

var routers = map[string]*chi.Mux{}
var cfgFile string

func globalHandler(w http.ResponseWriter, r *http.Request) {
	uriParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")

	prefix := uriParts[0]
	uriParts = uriParts[1:]
	newPath := "/" + strings.Join(uriParts, "/")

	r.URL.Path = newPath // Reset path to remove prefix
	// r.RequestURI = newPath // Probably not needed

	router, ok := routers[prefix]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No router found for %s", prefix)
		return
	}

	router.ServeHTTP(w, r)
}

func makeRouters() {
	if cfgFile == "" {
		panic("No config file specified")
	}

	configData, err := os.ReadFile(cfgFile)

	if err != nil {
		panic(err)
	}

	var config []types.TargetConfig

	json.Unmarshal(configData, &config)

	for _, target := range config {
		// Fallback to proxy
		setupRouter(target)
	}
}

func setupRouter(target types.TargetConfig) {
	router := chi.NewRouter()
	// router.Use(middleware.Logger)

	proxy := makeProxy(target.Target)

	for _, endpoint := range target.Endpoints {
		var method func(pattern string, handlerFn http.HandlerFunc)

		switch endpoint.Method {
		case "GET":
			method = router.Get
		case "POST":
			method = router.Post
		case "PUT":
			method = router.Put
		case "DELETE":
			method = router.Delete
		default:
			fmt.Println("Unknown method", endpoint.Path, endpoint.Method)
			continue
		}

		method(endpoint.Path, func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		})
	}

	if target.RefuseFallback == false {
		router.NotFound(proxy.ServeHTTP) // Catch all router
	}

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
	}
}

// gatewayCmd represents the gateway command
var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Start the gateway",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		makeRouters()

		globalHandler := http.HandlerFunc(globalHandler)

		router := chi.NewRouter()
		router.Use(middleware.Logger)

		router.NotFound(globalHandler)

		err := http.ListenAndServe("127.0.0.1:8000", router)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gatewayCmd)

	gatewayCmd.Flags().StringVarP(&cfgFile, "config", "c", "config.json", "Config File")
}
