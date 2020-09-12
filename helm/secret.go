package main

import (
	"context"
	"log"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateKeycloakSecret ...
func CreateKeycloakSecret(name, namespace, username, password string, realm []byte) (err error) {

	deploySecret := Clientset.CoreV1().Secrets(namespace)

	secret := &apiv1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "apps/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: map[string][]byte{
			"KEYCLOAK_USER":     []byte(username),
			"KEYCLOAK_PASSWORD": []byte(password),
			"realm.json":        realm,
		},
		Type: "Opaque",
	}

	result, err := deploySecret.Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	log.Printf("secret created: %q\n", result.GetObjectMeta().GetName())
	return
}

// DeleteSecret ...
func DeleteSecret(name, namespace string) (err error) {

	deploySecret := Clientset.CoreV1().Secrets(namespace)

	deletePolicy := metav1.DeletePropagationForeground
	if err := deploySecret.Delete(context.TODO(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	log.Println("secret deleted")
	return
}
