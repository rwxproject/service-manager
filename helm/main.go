package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/client-go/kubernetes"
)

// Clientset ...
var (
	Clientset    *kubernetes.Clientset
	ActionConfig *action.Configuration
)

func init() {
	gotenv.Load()
	Clientset, ActionConfig = Kubeinit()
}

func main() {
	log.Printf("version: %s", os.Getenv("VERSION"))
	r := mux.NewRouter()
	r.HandleFunc("/deploy", DeployHandler).Methods("POST")
	r.HandleFunc("/status", StatusHandler).Methods("GET")
	r.HandleFunc("/list", ListHandler).Methods("GET")
	log.Println("server running on port 8899")
	log.Fatal(http.ListenAndServe(":8899", r))
}
