package main

import (
	"encoding/json"
	"log"
	"net/http"

	"helm.sh/helm/v3/pkg/action"
	helmtime "helm.sh/helm/v3/pkg/time"
)

// ReleaseInfo ...
type ReleaseInfo struct {
	Name         string        `json:"name"`
	Namespace    string        `json:"namespace"`
	Revision     int           `json:"revision"`
	Updated      helmtime.Time `json:"updated"`
	Status       string        `json:"status"`
	ChartName    string        `json:"chartName"`
	ChartVersion string        `json:"chartVersion"`
	AppVersion   string        `json:"appVersion"`
	Description  string        `json:"description"`
}

// ListHandler ..
func ListHandler(w http.ResponseWriter, r *http.Request) {
	res, err := ListStatus("default")
	if err != nil {
		log.Printf(err.Error())
	}
	json.NewEncoder(w).Encode(res)
}

// ListStatus ..
func ListStatus(namespace string) (res []ReleaseInfo, err error) {

	client := action.NewList(ActionConfig)
	client.Deployed = true

	results, err := client.Run()
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	for _, rel := range results {
		var r ReleaseInfo
		r.Name = rel.Name
		r.Namespace = rel.Namespace
		r.Revision = rel.Version
		r.Updated = rel.Info.LastDeployed
		r.Status = string(rel.Info.Status)
		r.ChartName = rel.Chart.Name()
		r.ChartVersion = rel.Chart.Metadata.Version
		r.AppVersion = rel.Chart.AppVersion()
		r.Description = rel.Info.Description
		res = append(res, r)
	}
	return res, nil
}
