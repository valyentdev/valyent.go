_This repo is now archived._

# Valyent Go SDK

## Installation

```bash
go get github.com/valyentdev/valyent.go
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/valyentdev/valyent.go"
    "github.com/valyentdev/ravel/api"
)

func main() {
    // Initialize the client
    client := valyent.NewClient().
        WithBearerToken("your-api-token")

    // Create a new fleet
    fleet, err := client.CreateFleet(api.CreateFleetPayload{
        Name: "my-fleet",
    })
    if err != nil {
        panic(err)
    }

    // Create a machine in the fleet
    machine, err := client.CreateMachine(fleet.ID, api.CreateMachinePayload{
        Region: "us-east-1",
        Config: api.MachineConfig{
            Image: "nginx:latest",
            Guest: api.GuestConfig{
                CPUKind:  "shared",
                MemoryMB: 512,
                CPUs:     1,
            },
            Workload: api.Workload{
                Env: []string{"PORT=8080"},
            },
        },
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Created machine: %s\n", machine.ID)
}
```

## Core Concepts

### Client

The main entry point for the SDK is the `Client` type. Create a new client and configure it with your API token:

```go
client := valyent.NewClient().
    WithBearerToken("your-api-token").
    WithBaseURL("https://console.valyent.cloud") // Optional, defaults to this value
```

### Fleets

```go
// Create a fleet
fleet, err := client.CreateFleet(api.CreateFleetPayload{
    Name: "production",
})

// List all fleets
fleets, err := client.GetFleets()

// Delete a fleet
err := client.DeleteFleet(fleetID)
```

### Machines

```go
// Create a machine
machine, err := client.CreateMachine(fleetID, api.CreateMachinePayload{
    Region: "us-east-1",
    Config: api.MachineConfig{
        Image: "my-image:latest",
        Guest: api.GuestConfig{
            CPUKind:  "shared",
            MemoryMB: 1024,
            CPUs:     2,
        },
    },
})

// List machines in a fleet
machines, err := client.GetMachines(fleetID)

// Get machine events
events, err := client.GetMachineEvents(fleetID, machineID)

// Start a machine
err := client.StartMachine(fleetID, machineID)

// Stop a machine
err := client.StopMachine(fleetID, machineID)

// Delete a machine
err := client.DeleteMachine(fleetID, machineID, false)
```

### Gateways

```go
// Create a gateway
gateway, err := client.CreateGateway(api.CreateGatewayPayload{
    Name:       "web-gateway",
    TargetPort: 8080,
})

// List gateways
gateways, err := client.GetGateways()

// Delete a gateway
err := client.DeleteGateway(gatewayID)
```

### Logs and Events

```go
// Get machine logs
logs, err := client.GetLogs(fleetID, machineID)

// Stream logs in real-time
ctx := context.Background()
stream, err := client.StreamLogs(ctx, valyent.LogStreamOptions{
    FleetID:   fleetID,
    MachineID: machineID,
})
defer stream.Close()

for {
    entry, ok := stream.Next()
    if !ok {
        if err := stream.Err(); err != nil {
            panic(err)
        }
        break
    }
    fmt.Printf("[%s] %s\n", entry.Level, entry.Message)
}
```

### Environment Variables

```go
// Get environment variables
env, err := client.GetEnvironmentVariables(namespace, fleetID)

// Set environment variables
redeploy, err := client.SetEnvironmentVariables(namespace, fleetID, []string{
    "KEY1=value1",
    "KEY2=value2",
})
```

### Deployments

```go
// Create a deployment with a tarball
deployment, err := client.CreateDeployment(
    namespace,
    fleetID,
    api.CreateDeploymentPayload{
        Machine: machineConfig,
    },
    tarballReader,
)
```

## Error Handling

The SDK provides detailed error information through Go's standard error interface:

```go
machine, err := client.CreateMachine(fleetID, config)
if err != nil {
    fmt.Printf("Failed to create machine: %v\n", err)
    return
}
```

## Best Practices

1. **Context Usage**: Use context for operations that might need cancellation or timeouts, especially with streaming logs:

   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()

   stream, err := client.StreamLogs(ctx, options)
   ```

2. **Resource Cleanup**: Always close resources that implement `io.Closer`:

   ```go
   stream, err := client.StreamLogs(ctx, options)
   if err != nil {
    return err
   }
   defer stream.Close()
   ```

3. **Error Checking**: Always check errors returned by SDK methods:
   ```go
   if err := client.DeleteMachine(fleetID, machineID, false); err != nil {
       log.Printf("Failed to delete machine: %v", err)
       // Handle error appropriately
   }
   ```
