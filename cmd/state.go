package cmd

import (
	"log"

	"github.com/mythicmc/peridot/config"
	"github.com/mythicmc/peridot/repos"
)

func HandleStateCommand(args []string) {
	// Load repositories
	repositories, err := repos.LoadRepositories()
	if err != nil {
		log.Fatalln("An error has occurred while loading repositories:", err)
		return
	}
	// Load configuration
	_, err = config.LoadConfigs(repositories)
	if err != nil {
		log.Fatalln("An error has occurred while loading configuration:", err)
		return
	}

	// FIXME: Prepare changes to software
	// FIXME: Prepare changes to plugins
	// FIXME: Prepare changes to server.properties

	// FIXME: Display changes to software
	// FIXME: Display changes to plugins
	// FIXME: Display changes to server.properties
}
