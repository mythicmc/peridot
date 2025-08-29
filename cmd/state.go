package cmd

import (
	"fmt"
	"log"

	"github.com/mythicmc/peridot/config"
	"github.com/mythicmc/peridot/deploy"
	"github.com/mythicmc/peridot/repos"
	"github.com/mythicmc/peridot/utils"
)

func loadReposConfigUpdateState() (
	repos.Repositories,
	config.Configs,
	map[string]deploy.SoftwareUpdateOperation,
	map[string][]deploy.ServerPropertiesUpdateOperation,
	map[string]map[string]deploy.PluginUpdateOperation,
) {
	// Load repositories
	repositories, err := repos.LoadRepositories()
	if err != nil {
		log.Fatalln("An error has occurred while loading repositories:", err)
	}
	// Load configuration
	configs, err := config.LoadConfigs(repositories)
	if err != nil {
		log.Fatalln("An error has occurred while loading configuration:", err)
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

	return repositories, configs, softwareUpdates, serverPropertiesUpdates, pluginUpdates
}

func previewUpdates(
	softwareUpdates map[string]deploy.SoftwareUpdateOperation,
	serverPropertiesUpdates map[string][]deploy.ServerPropertiesUpdateOperation,
	pluginUpdates map[string]map[string]deploy.PluginUpdateOperation,
) {
	if len(softwareUpdates) > 0 {
		fmt.Println("Pending software updates:")
		for server, update := range softwareUpdates {
			prevHash := utils.PickNonEmptyString(update.PrevHash, "(missing)")
			newHash := utils.PickNonEmptyString(update.NewHash, "(removed)")
			fmt.Printf(" - %s: %s (%s -> %s)\n", server, update.SoftwareType, prevHash, newHash)
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
				oldValue := utils.PickNonEmptyString(update.OldValue, "(missing)")
				newValue := utils.PickNonEmptyString(update.NewValue, "(removed)")
				fmt.Printf("\t=> %s: %s -> %s\n", update.Property, oldValue, newValue)
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
				prevVersion := utils.PickNonEmptyString(update.PrevVersion, "(missing)")
				newVersion := utils.PickNonEmptyString(update.NewVersion, "(removed)")
				fmt.Printf("    => %s: %s -> %s\n", update.PluginName, prevVersion, newVersion)
			}
		}
	} else {
		fmt.Println("Plugins: Up to date")
	}
}

func HandleStateCommand(args []string) {
	_, _, softwareUpdates, serverPropertiesUpdates, pluginUpdates := loadReposConfigUpdateState()
	previewUpdates(softwareUpdates, serverPropertiesUpdates, pluginUpdates)
}
