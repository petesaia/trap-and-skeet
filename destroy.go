package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func destroy(ctx context.Context, cli *client.Client, id string, done chan bool) error {
	if err := cli.ContainerStop(ctx, id, nil); err != nil {
		return err
	}
	done <- true
	return nil
}

// Remove all packages created with Trap and Skeet. This can remain simple for
// now. More configs may be added later.
func Remove() {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	done := make(chan bool)

	// Pull a list of all running containers.
	containers, err := cli.ContainerList(
		ctx,
		types.ContainerListOptions{},
	)
	if err != nil {
		panic(err)
	}

	// Grab the number of containers.
	numOfContainers := len(containers)

	for _, container := range containers {
		if strings.Contains(container.Names[0], Identifier) {
			go destroy(ctx, cli, container.ID, done)
		} else {
			numOfContainers--
		}
	}

	if numOfContainers == 1 {
		Log("Destroying 1 container.")
	} else if numOfContainers > 1 {
		Log("Destroying " + fmt.Sprintf("%d", numOfContainers) + " containers.")
	} else {
		Log("Nothing to destroy!")
	}

	s.Start()

	for i := 0; i < numOfContainers; i++ {
		<-done
	}

	s.Stop()
}
