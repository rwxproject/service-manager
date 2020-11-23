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
	Ready  bool   `json:"ready"`
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
				st.Ready = string(s.Status) == "True"
			}
		}
		log.Printf("%v: %+v", pod.GetObjectMeta().GetName(), pod.Status.Phase)
		res = append(res, st)
	}
	return res, nil
}

// SvcStatus ...
type SvcStatus struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// SvcHandler ..
func SvcHandler(w http.ResponseWriter, r *http.Request) {
	res, err := SvcsStatus("default")
	if err != nil {
		log.Printf(err.Error())
	}
	json.NewEncoder(w).Encode(res)
}

// SvcsStatus ...
func SvcsStatus(namespace string) (p []SvcStatus, err error) {
	svcs, err := Clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	var res []SvcStatus
	for _, svc := range svcs.Items {
		var st SvcStatus
		st.Name = svc.GetObjectMeta().GetName()
		st.Type = string(svc.Spec.Type)
		log.Printf("%v: %+v", svc.GetObjectMeta().GetName(), svc.Spec.Type)
		res = append(res, st)
	}
	return res, nil
}
