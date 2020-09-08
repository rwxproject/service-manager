package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// DeployPlayload ...
type DeployPlayload struct {
	Namespace string `json:"namespace"`
	Keycloak  struct {
		Realm string `json:"realm"`
		User  struct {
			Username  string `json:"username"`
			Password  string `json:"password"`
			Privilege bool   `json:"privilege"`
		} `json:"user"`
	} `json:"keycloak"`
}

// DeployHandler ..
func DeployHandler(w http.ResponseWriter, r *http.Request) {
	var deployPlayload DeployPlayload

	err := json.NewDecoder(r.Body).Decode(&deployPlayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Panic(err.Error())
	}

	// keycloak realm template
	tpl, err := template.ParseFiles("realm-tpl.json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Fatalln(err)
	}

	// create keycloak realm file
	rm, err := os.Create("realm.json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Fatalln("error creating file", err)
	}
	defer rm.Close()

	// write keycloak realm file
	err = tpl.Execute(rm, deployPlayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Fatalln(err)
	}

	// read keycloak realm file
	jsonFile, err := os.Open("realm.json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		fmt.Println(err)
	}
	defer jsonFile.Close()

	// get realm.json
	realm, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		fmt.Println(err)
	}

	// delete old secret if exists
	DeleteSecret()
	time.Sleep(2 * time.Second)

	// create secret
	errs := CreateSecret(deployPlayload.Keycloak.User.Username, deployPlayload.Keycloak.User.Password, realm)
	if errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errs.Error()))
		return
	}

	// install helm chart
	chartPath := "/Users/fabiano/rwxproject/kube-manager/chart/keycloak-9.0.8.tgz"
	erri := Install("keycloak", deployPlayload.Namespace, chartPath)
	if erri != nil {
		DeleteSecret()
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(erri.Error()))
		return
	}

	// service status
	res, err := ServiceStatus("default")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(res)

}
