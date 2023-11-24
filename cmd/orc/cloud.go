package main

import (
	"errors"
	"fmt"
	v1 "github.com/gophercloud/openstack-resource-controller/api/v1alpha1"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

func NewOpenStackCloud(cloudName string) v1.OpenStackCloud {
	return v1.OpenStackCloud{
		Spec: v1.OpenStackCloudSpec{
			Cloud: cloudName,
			Credentials: v1.OpenStackCloudCredentials{
				Source: "secret",
				SecretRef: v1.OpenStackCloudCredentialsSecretRef{
					Name: "openstackcloud-" + cloudName,
					Key:  "clouds.yaml",
				},
			},
		},
	}
}

func cloud(_ []string) error {
	cloudName := coalesce(opts.OsCloud, os.Getenv("OS_CLOUD"))
	if cloudName == "" {
		return fmt.Errorf("OS_CLOUD environment variable not found")
	}

	locations := []string{"clouds.yaml"}
	if configDir, err := os.UserConfigDir(); err == nil {
		locations = append(locations, path.Join(configDir, "openstack", "clouds.yaml"))
	}

	for _, location := range locations {
		_, err := os.Stat(location)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return fmt.Errorf("unable to stat %q: %v", location, err)
		}

		f, err := os.Open(location)
		if err != nil {
			return fmt.Errorf("unable to open %q: %v", location, err)
		}

		var clouds struct {
			Clouds map[string]interface{} `yaml:"clouds"`
		}
		if err := yaml.NewDecoder(f).Decode(&clouds); err != nil {
			return fmt.Errorf("unable to decode %q: %v", location, err)
		}

		cloud, ok := clouds.Clouds[cloudName]
		if !ok {
			return fmt.Errorf("cloud %q not found in %q", cloudName, location)
		}

		clouds.Clouds = map[string]interface{}{cloudName: cloud}

		// return yaml.Marshal(clouds)
		b, err := yaml.Marshal(NewOpenStackCloud(cloudName))
		if err != nil {
			return err
		}
		fmt.Printf("%s", b)
		return nil
	}

	return fmt.Errorf("Error: could not find clouds.yaml in the default locations.")
}
