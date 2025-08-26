package cmd

import (
	"fmt"
	"log"

	"github.com/mythicmc/peridot/config"
	"github.com/mythicmc/peridot/deploy"
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
	configs, err := config.LoadConfigs(repositories)
	if err != nil {
		log.Fatalln("An error has occurred while loading configuration:", err)
		return
	}

	// Prepare changes to software
	softwareUpdates, err := deploy.PrepareAllSoftwareUpdates(repositories, configs)
	if err != nil {
		log.Println("An error has occurred while preparing software updates:", err,
			"Continuing without software updates...")
	}
	// Prepare changes to server.properties
	serverPropertiesUpdates, err := deploy.PrepareAllServerPropertiesUpdates(configs)
	if err != nil {
		log.Println("An error has occurred while preparing server properties updates:", err,
			"Continuing without server properties updates...")
	}
	// Prepare changes to plugins
	pluginUpdates, err := deploy.PrepareAllPluginUpdates(repositories, configs)
	if err != nil {
		log.Println("An error has occurred while preparing plugin updates:", err,
			"Continuing without plugin updates...")
	}

	if len(softwareUpdates) > 0 {
		fmt.Println("Pending software updates:")
		for server, update := range softwareUpdates {
			fmt.Printf(" - %s: %s (%s -> %s)\n", server, update.SoftwareType, update.PrevHash, update.NewHash)
		}
	} else {
		fmt.Println("Software: Up to date")
	}
	fmt.Println("==============================")
	if len(serverPropertiesUpdates) > 0 {
		fmt.Println("Pending server properties updates:")
		for server, updates := range serverPropertiesUpdates {
			fmt.Println(" - " + server)
			for _, update := range updates {
				fmt.Printf("\t=> %s: %s -> %s\n", update.Property, update.OldValue, update.NewValue)
			}
		}
	} else {
		fmt.Println("Server properties: Up to date")
	}
	fmt.Println("==============================")
	if len(pluginUpdates) > 0 {
		fmt.Println("Pending plugin updates:")
		for server, updates := range pluginUpdates {
			fmt.Println(" -> " + server)
			for _, update := range updates {
				fmt.Printf("    => %s: %s -> %s\n", update.PluginName, update.PrevVersion, update.NewVersion)
			}
		}
	} else {
		fmt.Println("Plugins: Up to date")
	}
}
