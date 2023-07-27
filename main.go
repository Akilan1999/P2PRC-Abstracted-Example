package main

import (
    "fmt"
    "github.com/Akilan1999/p2p-rendering-computation/abstractions"
    "github.com/Akilan1999/p2p-rendering-computation/client"
    "github.com/Akilan1999/p2p-rendering-computation/config"
    "github.com/Akilan1999/p2p-rendering-computation/p2p"
    "github.com/Akilan1999/p2p-rendering-computation/server/docker"
    "os"
    "time"
)

func main() {
    var Config *config.Config
    // check if the config file exists
    if _, err := os.Stat("config.json"); err != nil {
        // Initialize with base p2prc config files
        Config, err = abstractions.Init(nil)
        if err != nil {
            fmt.Println(err)
            return
        }
    }

    fmt.Println(Config)

    // start p2prc
    _, err := abstractions.Start()
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("Sleeping for 2 seconds")
    time.Sleep(time.Second * 2)

    table, err := p2p.ReadIpTable()
    if err != nil {
        fmt.Println(err)
        return
    }

    p2p.PrintIpTable()

    // Iterate through all nodes available
    // and spawn based on certain parametes
    for _, address := range table.IpAddress {
        MainAddress := address.Ipv4 + ":" + address.ServerPort
        specs, err := client.GetSpecs(MainAddress)
        if err != nil {
            fmt.Println(err)
            return
        }

        if specs.RAM > 5000 {
            fmt.Println("Above 5GB")
        }
        fmt.Println(specs)
    }

    container, err := StartDockerContainer("0.0.0.0:8088", 1)
    if err != nil {
        fmt.Println(err)
    }

    err = RemoveDockerContainer("0.0.0.0:8088", container.ID)
    if err != nil {
        fmt.Println(err)
    }
}

// StartDockerContainer Starts docker container
func StartDockerContainer(IPAddress string, NumberOfPorts int) (*docker.DockerVM, error) {
    // Creates container and returns-back result to
    // access container
    resp, err := client.StartContainer(IPAddress, NumberOfPorts, false, "")
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// RemoveDockerContainer Removes Docker Container
func RemoveDockerContainer(IPAddress string, ContainerID string) error {
    err := client.RemoveContianer(IPAddress, ContainerID)
    if err != nil {
        return err
    }
    return nil
}
