package main

import (
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
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

// type yamlConfig struct {
// 	Service
// }

// var config = `# config.yaml
// service:
//   type: ClusterIP
//   httpsPort: 8443
// `

// Install helm chart
func Install(name, namespace, chartPath string, opt DeployPlayload) (err error) {

	// var config = "{\"service\":{\"type\":\"LoadBalancer\",\"httpPort\":\"8086\",\"httpsPort\":\"8447\"}}"

	chart, err := loader.Load(chartPath)
	if err != nil {
		log.Panic(err)
		return err
	}

	client := action.NewInstall(ActionConfig)
	client.Namespace = namespace
	client.ReleaseName = name

	var options releaseOptions
	// httpPort := fmt.Sprintf("service.httpPort=%v", opt.HTTPPort)
	// httpsPort := fmt.Sprintf("service.httpsPort=%v", opt.HTTPSPort)
	// options.SetValues = []string{"service.httpPort=8081"}
	options.SetValues = []string{}
	options.SetValues = append(options.SetValues, fmt.Sprintf("service.httpPort=%v", opt.HTTPPort))
	options.SetValues = append(options.SetValues, fmt.Sprintf("service.httpsPort=%v", opt.HTTPSPort))
	options.SetStringValues = []string{}
	// options.Values = config

	vals, err := mergeValues(options)
	if err != nil {
		return err
	}
	fmt.Println(vals)

	rel, err := client.Run(chart, vals)
	if err != nil {
		fmt.Println("_+_+_+_+_+")
		log.Panic(err)
		return err
	}
	log.Println("release installed: ", rel.Name)
	return nil
}

func mergeValues(options releaseOptions) (map[string]interface{}, error) {
	spew.Dump(options)
	vals := map[string]interface{}{}

	err := yaml.Unmarshal([]byte(options.Values), &vals)
	if err != nil {
		return vals, fmt.Errorf("failed parsing values")
	}
	spew.Dump(vals)
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

	return vals, nil
}
