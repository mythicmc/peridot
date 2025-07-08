# peridot

IaC deployment engine for Minecraft servers. Build, configure and update Minecraft servers from a single config file!

## Commands

- `peridot status` / `peridot state`: Displays the current status of the configured Minecraft servers, comparing the current state of the server files with the desired state defined in Peridot's configuration.
- `peridot apply`: Applies the desired state, specified in Peridot's configuration and repository data, to the Minecraft server files on-disk.

## Architecture

Peridot is built around 3 main components:

- [Repository Loader](#repository-loader)
- [Config Loader](#config-loader)
- State Compiler (planned)

  This component will generate the final state of the Minecraft server based on the configuration and repository data in a deterministic manner. A rebuilt and simplified version of the Deployment Engine will then be fed the final state to apply to the Minecraft server.
- [Deployment Engine](#deployment-engine)

### Repository Loader

This component is responsible for managing Peridot repositories, located in the `repos` folder. These repositories contain Minecraft server JAR files and plugin files. It loads all data about on-disk repositories, performs validations (e.g. duplicate plugin files), then provides a list of available repositories to the Config Loader.

### Config Loader

This component is responsible for loading Peridot server configuration files, located in the `configs` folder. It loads all configuration files, performs validations (e.g. duplicate server names, missing plugins from the repositories, etc), then provides the list of available servers to the Deployment Engine.

The loader leverages a JavaScript-based configuration format, allowing for easy extensibility and customization. Basic support for `require()` and `module.export` is implemented (not full CommonJS!) to share configuration info between files.

### Deployment Engine

This component is responsible for deploying Minecraft servers based on the configuration and software provided by the Config/Repository Loaders. It checks the configuration and software against the on-disk server files, diffing the current state of the server with the desired state. Upon request, it can then apply the necessary changes to the server files, ensuring that the server is in the desired state.

It features integrations for Octyne, allowing Minecraft servers to be automatically restarted after applying the wanted state.
