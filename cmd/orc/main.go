// package main implements a subset of the commands available in
// python-openstackclient. Instead of directly targeting a cloud, it creates
// the corresponding openstack-resource-controller CRs.
package main

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

var Version = "dev"

var opts struct {
	OsCloud string `long:"os-cloud" description:"Cloud name in clouds.yaml"`
}

func main() {
	if len(os.Args) < 2 {
		return
	}

	switch resource, args := os.Args[1], os.Args[2:]; resource {
	case "help":
		log.Print("orc imitates the commands of the openstack cli and creates CRDs instead of OpenStack resources")
		return
	case "version":
		log.Print("github.com/gophercloud/openstack-resource-controller")
		log.Printf("Version %q", Version)
		return
	case "cloud":
		if err := cloud(args); err != nil {
			log.Fatalf(err.Error())
		}
	case "network":
		if err := network(args); err != nil {
			log.Fatalf(err.Error())
		}
	default:
		log.Fatalf("Fatal: unknown resource %q\n", resource)
	}
}

func network(args []string) error {
	panic("Not implemented yet!")
}

// coalesce returns the first non-empty argument, or the empty string.
func coalesce(arguments ...string) string {
	for i := range arguments {
		if arguments[i] != "" {
			return arguments[i]
		}
	}
	return ""
}

func init() {
	// Remove the timestamp prefix from the default logger
	log.SetFlags(0)

	if _, err := flags.Parse(&opts); err != nil {
		log.Fatalf("Unable to parse command-line arguments: %v", err)
	}
}
