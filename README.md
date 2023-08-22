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

1. **Build and Prepare**: Execute `make start-local` in your command line. This will build both the `sampleplugin` and the runtime, and place them in the appropriate folder structure.

2. **Launch Runtime**: After the build process is complete, run the following command to launch the runtime:

    ```bash
    ./bin/ar
    ```

### Container Environment

1. **Create Docker Images**: Run the following command to build Docker images for the runtime and the package repository:

    ```bash
    make start
    ```

2. **Launch Containers**: The above command will also automatically launch the created Docker containers for the runtime and package repository.

---

