package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ghodss/yaml"
	"gopkg.in/square/go-jose.v2/json"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"k8s.io/helm/pkg/strvals"
)

type releaseOptions struct {
	// common
	DryRun          bool          `json:"dry_run"`
	DisableHooks    bool          `json:"disable_hooks"`
	Wait            bool          `json:"wait"`
	Devel           bool          `json:"devel"`
	Description     string        `json:"description"`
	Atomic          bool          `json:"atomic"`
	SkipCRDs        bool          `json:"skip_crds"`
	SubNotes        bool          `json:"sub_notes"`
	Timeout         time.Duration `json:"timeout"`
	Values          string        `json:"values"`
	SetValues       []string      `json:"set"`
	SetStringValues []string      `json:"set_string"`

	// only install
	CreateNamespace  bool `json:"create_namespace"`
	DependencyUpdate bool `json:"dependency_update"`

	// only upgrade
	Force         bool `json:"force"`
	Install       bool `json:"install"`
	Recreate      bool `json:"recreate"`
	CleanupOnFail bool `json:"cleanup_on_fail"`
}

// ReleaseInstall helm chart
func ReleaseInstall(name, namespace, chartPath string, setValues []string) (err error) {

	chart, err := loader.Load(chartPath)
	if err != nil {
		log.Panic(err)
		return err
	}

	client := action.NewInstall(ActionConfig)
	client.Namespace = namespace
	client.ReleaseName = name
	// client.CreateNamespace = true
	var options releaseOptions
	options.SetStringValues = []string{}
	options.SetValues = setValues

	vals, err := mergeValues(options)
	if err != nil {
		return err
	}
	fmt.Println(vals)

	rel, err := client.Run(chart, vals)
	if err != nil {
		return err
	}
	log.Println("release installed: ", rel.Name)
	return nil
}

// ReleaseUninstall helm chart
func ReleaseUninstall(name, namespace string) (err error) {

	client := action.NewUninstall(ActionConfig)

	rel, err := client.Run(name)
	if err != nil {
		return err
	}
	log.Println("release uninstalled: ", rel.Release.Name)
	return nil
}

func mergeValues(options releaseOptions) (map[string]interface{}, error) {
	vals := map[string]interface{}{}
	//
	yamlFile, err := ioutil.ReadFile("values.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	val1, err1 := yaml.YAMLToJSON(yamlFile)

	// err := yaml.Unmarshal([]byte(options.Values), &vals)
	if err1 != nil {
		return vals, fmt.Errorf("failed parsing values")
	}
	err2 := json.Unmarshal(val1, &vals)
	if err2 != nil {
		return vals, fmt.Errorf("a failed parsing values")
	}
	for _, value := range options.SetValues {
		if err := strvals.ParseInto(value, vals); err != nil {
			return vals, fmt.Errorf("failed parsing set data")
		}
	}
	for _, value := range options.SetStringValues {
		if err := strvals.ParseIntoString(value, vals); err != nil {
			return vals, fmt.Errorf("failed parsing set_string data")
		}
	}

	spew.Dump(vals)
	return vals, nil
}

func convert(m map[interface{}]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m {
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			res[fmt.Sprint(k)] = convert(v2)
		default:
			res[fmt.Sprint(k)] = v
		}
	}
	return res
}
