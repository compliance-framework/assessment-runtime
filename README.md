# Cloud Compliance Framework - Assessment Runtime

This documentation provides a complete overview of the Compliance Assessment Runtime System, a project designed to manage and execute plugins for various compliance checks or assessments. It is implemented primarily in Golang and uses gRPC for communication between services.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
	- [Plugin Management](#plugin-management)
	- [Manager Functions](#manager-functions)
- [Configuration Management](#configuration-management)
- [Plugin Interface](#plugin-interface)
- [Communication and Registration](#communication-and-registration)
	- [gRPC Server](#grpc-server)
- [Plugin Download](#plugin-download)
- [Integration and Dependencies](#integration-and-dependencies)
- [Summary](#summary)

## Overview

The main application serves as a centralized hub for running compliance checks or assessments via various plugins. Written in Golang, the system leverages the gRPC protocol for robust inter-service communication.

## Architecture

### Plugin Management

The architecture employs a plugin-based design, making each plugin an independent entity that can be managed, scheduled, and executed separately. Plugins are executed in separate processes, supported by context and synchronization primitives.

### Plugin Manager

Plugin Manager, implemented in `manager.go`, controls several critical aspects:

- **Initialization**: Utilizes the plugin's interface implementation for initial setup.
- **Scheduling**: Uses a cron-like scheduler to execute plugins based on their defined schedules.
- **Execution**: Handles the actual running of plugins, ensuring synchronization and error handling.

## Configuration Management

- **File Format**: The configuration is read from a YAML file upon the application's launch.
- **Function**: Utilizes the `LoadConfig` function in `manager.go`.
- **Data Structure**: Configuration details are stored in a `Config` struct defined in `models.go`.
- **Details**: The configuration includes essential attributes like the control plane, plugin name, version, and schedule.

## Plugin Interface

- **Specification**: The required interface is defined in `plugin.go`.
- **Methods**: Each plugin must implement three methods: Initialization, Execution, and Shutdown.
- **Extensibility**: Additional methods may be required based on specific project requirements.

## Communication and Registration

### gRPC Server

- **Implementation**: Each plugin operates as a separate gRPC server, facilitating robust communication.
- **Registration**: The Manager uses `register.go` to register each plugin server, enhancing modularity and scalability.

## Plugin Download

- **Functionality**: Managed by `downloader.go`, this feature allows the system to download plugins from a specified remote registry.
- **Dynamic Updates**: The system can dynamically add or update plugins without requiring a restart or full rebuild.

## Integration and Dependencies

The system integrates seamlessly with various components and dependencies:

- Plugin execution is scheduled via `scheduler.go`.
- Plugins are downloaded from remote sources as indicated by `downloader.go`.
- The Manager, implemented in `manager.go`, orchestrates the complete lifecycle of plugin objects.

## Running the Compliance Assessment Runtime

This section explains how to run the Compliance Assessment Runtime System either on your local system or in a containerized environment.

### Table of Contents

- [Local Environment](#local-environment)
- [Container Environment](#container-environment)

### Local Environment

1. **Setup**: Before you run `make run-local`, you need to:

- Install nats server

- Configure nats-server with jetstream, eg

```
host: 127.0.0.1
port: 4222
jetstream: enabled
```

- Ensure go 1.21+ is installed

2. **Build and Prepare**: Execute `make run-local` in your command line. This will build both the `sampleplugin` and the runtime, and place them in the appropriate folder structure.

### Container Environment

1. **Create Docker Images**: Run the following command to build Docker images for the runtime and the package repository:

    ```bash
    make compose-up
    ```

2. **Launch Containers**: The above command will also automatically launch the created Docker containers for the runtime and package repository.

---

# Plugin Manager

The `PluginManager` is a struct that manages the lifecycle of plugins in the application. It is responsible for starting, executing, and stopping plugins.

## Structure

The `PluginManager` struct has two fields:

- `cfg`: This is the configuration of the application.
- `clients`: This is a map where the keys are plugin names and the values are pointers to `goplugin.Client` instances.

## Methods

### NewPluginManager(cfg config.Config) *PluginManager

This function creates a new instance of `PluginManager`. It takes a `config.Config` object as an argument and initializes the `clients` map.

### Start() error

This method starts all the plugins defined in the configuration. It first groups the plugins by their package name. Then, for each package, it creates a new `goplugin.Client`, starts it, and stores it in the `clients` map.

### Execute(name string, input ActionInput) error

This method executes a specific plugin. It takes the name of the plugin and an `ActionInput` object as arguments. It finds the client for the plugin in the `clients` map, dispenses the plugin, and then executes it. If any step fails, it logs the error and returns it.

### Stop()

This method stops all the clients in the `clients` map. It does this concurrently using a `sync.WaitGroup`.

## Usage

To use the `PluginManager`, you first need to create a new instance with `NewPluginManager`, passing in your configuration. You can then start all the plugins with `Start`. To execute a specific plugin, use `Execute`, passing in the name of the plugin and an `ActionInput` object. Finally, you can stop all the plugins with `Stop`.

### Start() error

The `Start` method is responsible for initializing and starting all the plugins defined in the configuration.

Here's a step-by-step breakdown of what it does:

1. It creates a map called `pluginMap` where the keys are package names and the values are slices of `config.PluginConfig` objects.

2. It iterates over all the plugins in the configuration, grouping them by their package name in the `pluginMap`.

3. For each package in the `pluginMap`, it logs that it's loading the plugins for that package.

4. It then creates another map, also called `pluginMap`, where the keys are plugin names and the values are pointers to `AssessmentActionGRPCPlugin` instances.

5. It constructs the path to the plugin package using the package name and version from the first plugin in the package.

6. It logs that it's loading the plugin package, including the package name and path.

7. It creates a new `goplugin.Client` with the `HandshakeConfig`, the `pluginMap`, and a command to execute the plugin package. It also specifies that the client should use the GRPC protocol.

8. Finally, it stores the client in the `clients` map, using the package name as the key.

If any errors occur during this process, the method will return the error.
