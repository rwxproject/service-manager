package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodStatus ...
type PodStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Ready  string `json:"ready"`
}

// StatusHandler ..
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	res, err := PodsStatus("default")
	if err != nil {
		log.Printf(err.Error())
	}
	json.NewEncoder(w).Encode(res)
}

// PodsStatus ...
func PodsStatus(namespace string) (p []PodStatus, err error) {
	pods, err := Clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	var res []PodStatus
	for _, pod := range pods.Items {
		var st PodStatus
		st.Name = pod.GetObjectMeta().GetName()
		st.Status = string(pod.Status.Phase)
		for _, s := range pod.Status.Conditions {
			if s.Type == "Ready" {
				st.Ready = string(s.Status)
			}
		}
		log.Printf("%v: %+v", pod.GetObjectMeta().GetName(), pod.Status.Phase)
		res = append(res, st)
	}
	return res, nil
}
