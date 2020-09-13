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
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Chart     string `json:"chart"`
	HTTPPort  string `json:"httpPort"`
	HTTPSPort string `json:"httpsPort"`
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
	DeleteSecret("keycloak-setup", deployPlayload.Namespace)
	time.Sleep(1 * time.Second)

	// create secret
	errs := CreateKeycloakSecret("keycloak-setup", deployPlayload.Namespace, deployPlayload.Keycloak.User.Username, deployPlayload.Keycloak.User.Password, realm)
	if errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errs.Error()))
		return
	}

	// release install
	// chartPath := "/../chart/keycloak-9.0.8.tgz"
	chartPath := fmt.Sprintf("%s/%s.tgz", os.Getenv("CHART_PATH"), deployPlayload.Chart)
	fmt.Println(chartPath)

	var SetValues = []string{}
	SetValues = append(SetValues, fmt.Sprintf("service.httpPort=%v", deployPlayload.HTTPPort))
	SetValues = append(SetValues, fmt.Sprintf("service.httpsPort=%v", deployPlayload.HTTPSPort))

	erri := ReleaseInstall(deployPlayload.Name, deployPlayload.Namespace, chartPath, SetValues)
	if erri != nil {
		DeleteSecret("keycloak-setup", deployPlayload.Namespace)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(erri.Error()))
		return
	}
	// delete created secret
	time.Sleep(1 * time.Second)
	DeleteSecret("keycloak-setup", deployPlayload.Namespace)

	// helm ls
	res, err := ListStatus(deployPlayload.Namespace)
	if err != nil {
		log.Printf(err.Error())
	}
	json.NewEncoder(w).Encode(res)

}

// DeployUninstallHandler ..
func DeployUninstallHandler(w http.ResponseWriter, r *http.Request) {
	var deployPlayload DeployPlayload

	err := json.NewDecoder(r.Body).Decode(&deployPlayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Panic(err.Error())
	}

	// release uninstall
	erri := ReleaseUninstall(deployPlayload.Name, deployPlayload.Namespace)
	if erri != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(erri.Error()))
		return
	}

	// helm ls
	res, err := ListStatus(deployPlayload.Namespace)
	if err != nil {
		log.Printf(err.Error())
	}
	json.NewEncoder(w).Encode(res)

}
