package main

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SvcStatus ...
type SvcStatus struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// StatusHandler ..
// func StatusHandler(w http.ResponseWriter, r *http.Request) {
// 	res, err := SvcStatus("default")
// 	if err != nil {
// 		log.Printf(err.Error())
// 	}
// 	json.NewEncoder(w).Encode(res)
// }

// ServiceStatus ...
func ServiceStatus(namespace string) (p []SvcStatus, err error) {
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
