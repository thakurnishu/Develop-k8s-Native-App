package main

import (
	"flag"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func main() {

	resource := flag.String("res", "", "name of the resource")
	flag.Parse()

	configFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()

	matchVersionFlag := cmdutil.NewMatchVersionFlags(configFlags)

	mapper, err := cmdutil.NewFactory(matchVersionFlag).ToRESTMapper()
	if err != nil {
		fmt.Printf("getting restmapper from newfactory: %s\n", err.Error())
		return
	}

	gvr, err := mapper.ResourceFor(
		schema.GroupVersionResource{
			Resource: *resource,
		},
	)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("Complete GVR: group %s, version %s, resource %s\n", gvr.Group, gvr.Version, gvr.Resource)

}
