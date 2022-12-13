package commands

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"monot/config"
	"os"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

// InitCmd -
func InitCmd(cCtx *cli.Context) error {
	conf := &config.Config{}

	f, err := os.Open("./monot-manifest.yaml")
	if err != nil {
		log.Panic(err)
	}

	manifest, err := io.ReadAll(f)
	if err != nil {
		log.Panic(err)
	}

	if err := yaml.Unmarshal(manifest, &conf); err != nil {
		log.Panic(err)
	}

	conf.RegisteredServices = make(map[string]*config.Service)
	for _, s := range conf.Services {
		var service *config.Service
		f, err := os.Open(s + "/monot-service.yaml")
		if err != nil {
			log.Panic(err)
		}

		serviceDefinition, err := io.ReadAll(f)
		if err != nil {
			log.Panic(err)
		}

		if err := yaml.Unmarshal(serviceDefinition, &service); err != nil {
			log.Panic(err)
		}

		log.Println(service)

		conf.RegisteredServices[s] = service
	}

	if err := os.MkdirAll("./.monot/", 0777); err != nil {
		log.Panic(err)
	}

	output := &bytes.Buffer{}
	g := gob.NewEncoder(output)
	if err := g.Encode(conf); err != nil {
		log.Panic(err)
	}

	if err := os.WriteFile("./.monot/cache", output.Bytes(), 0777); err != nil {
		log.Panic(err)
	}

	fmt.Println("initialisation complete")

	return nil
}
