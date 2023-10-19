# Mikrotik Package

This Go package provides functionality to interact with Mikrotik devices using the RouterOS API.

## Installation

To use this package, you need to have Go installed. You can install the package using the following command:

```shell
go get github.com/OzkrOssa/mikrotik-go
```

## Usage

```go
package main

import (
	"fmt"
	"log"
	"github.com/OzkrOssa/mikrotik-go"
)

func main() {
	// Replace the following variables with your Mikrotik credentials
	address := "your-mikrotik-address"
	username := "your-username"
	password := "your-password"

	// Create a new Mikrotik instance
	mikrotikClient, err := mikrotik.NewMikrotikRepository(address, username, password)
	if err != nil {
		log.Fatalf("Failed to connect to Mikrotik: %v", err)
	}

	// Get the Mikrotik identity
	identity, err := mikrotikClient.GetIdentity()
	if err != nil {
		log.Fatalf("Failed to get Mikrotik identity: %v", err)
	}
	fmt.Println("Mikrotik Identity:", identity)

	// Get secrets and active connections
	secrets, err := mikrotikClient.GetSecrets("your-bts-name", "your-bts-host")
	if err != nil {
		log.Fatalf("Failed to get secrets: %v", err)
	}
	fmt.Println("Secrets:", secrets)

	activeConnections, err := mikrotikClient.GetActiveConnections()
	if err != nil {
		log.Fatalf("Failed to get active connections: %v", err)
	}
	fmt.Println("Active Connections:", activeConnections)

	// Enable SNMP on the Mikrotik
	mikrotikClient.EnableSNMP()

	// Set MAC address from active connections to secrets
	mikrotikClient.SetMacFromAC()

	// Set remote address from active connections to secrets
	mikrotikClient.SetRemoteAddress()

	// Add a secret to the address list
	ip := "192.168.1.100"
	comment := "Example Comment"
	addressList := "your-address-list"
	mikrotikClient.AddSecretToAddressList(ip, comment, addressList)
}
```

## Documentation

### NewMikrotikRepository

```go
func NewMikrotikRepository(addr, user, password string) (Mikrotik, error)
```

NewMikrotikRepository creates a new Mikrotik client instance and connects to the specified Mikrotik device using the given address, username, and password. It returns a Mikrotik interface and an error if any.

### Mikrotik Methods

- `GetIdentity() (map[string]string, error)`: Retrieves the identity information of the Mikrotik device.

- `GetSecrets(btsName string, btsHost string) ([]map[string]string, error)`: Retrieves secrets from the Mikrotik device. The `btsName` and `btsHost` parameters are used to add additional information to the secrets.

- `GetActiveConnections() ([]map[string]string, error)`: Retrieves active connections from the Mikrotik device.

- `EnableSNMP()`: Enables SNMP on the Mikrotik device.

- `SetMacFromAC()`: Sets MAC addresses from active connections to secrets.

- `SetRemoteAddress()`: Sets remote addresses from active connections to secrets.

- `AddSecretToAddressList(ip string, comment string, addressList string)`: Adds a secret to the specified address list.

