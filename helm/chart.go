package main

import (
	"log"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
)

// Install helm chart
func Install(name, namespace, chartPath string) (err error) {

	chart, err := loader.Load(chartPath)
	if err != nil {
		log.Panic(err)
		return err
	}

	client := action.NewInstall(ActionConfig)
	client.Namespace = namespace
	client.ReleaseName = name

	rel, err := client.Run(chart, nil)
	if err != nil {
		log.Panic(err)
		return err
	}
	log.Println("release installed: ", rel.Name)
	return nil
}
