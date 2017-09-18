package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// ConfigureInput is what's passed to spawn. It contains instructions for
// building the ephemeral containers.
type ConfigureInput struct {
	Image      string
	Port       string
	User       string
	Pass       string
	SSHKeyFile string
}

// Pull the image.
func pullImage(ctx context.Context, cli *client.Client, config *ConfigureInput) error {
	Log("Pulling docker image: " + config.Image)
	pullReader, err := cli.ImagePull(
		ctx,
		config.Image,
		types.ImagePullOptions{
			All: true,
		},
	)
	if err != nil {
		return err
	}
	defer pullReader.Close()
	if _, err = ioutil.ReadAll(pullReader); err != nil {
		return err
	}

	return nil
}

// Create a new container.
func createContainer(ctx context.Context, cli *client.Client, config *ConfigureInput) (*container.ContainerCreateCreatedBody, error) {
	Log("Creating container from image.")
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: config.Image,
			Cmd:   []string{"tail", "-f", "/dev/null"},
			ExposedPorts: nat.PortSet{
				"22/tcp": struct{}{},
			},
		},
		&container.HostConfig{
			AutoRemove: true,
			PortBindings: nat.PortMap{
				"22/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: config.Port,
					},
				},
			},
		},
		nil,
		"",
	)
	if err != nil {
		return nil, err
	}

	// Run the container.
	if err := cli.ContainerStart(
		ctx,
		resp.ID,
		types.ContainerStartOptions{},
	); err != nil {
		return nil, err
	}

	// Give the container an identifier.
	if err := cli.ContainerRename(
		ctx,
		resp.ID,
		RandomString(5)+"_"+Identifier+"_"+RandomString(5),
	); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Run all the command necessary on the poor container.
func runCommands(ctx context.Context, cli *client.Client, config *ConfigureInput, containerID string) error {
	if config.User != "root" {
		if err := cmd(
			ctx,
			cli,
			[]string{
				"adduser",
				"--disabled-password",
				"--gecos",
				"TaS User,666,202-237-4074,202-237-4074",
				config.User,
			},
			containerID,
		); err != nil {
			return err
		}
	}
	if config.Pass != "" {
		if err := cmd(
			ctx,
			cli,
			[]string{
				"sh",
				"-c",
				"printf 'PasswordAuthentication yes' >> /etc/ssh/sshd_config",
				"&&",
				"service",
				"ssh",
				"restart",
				"&&",
				"usermod",
				"--password",
				"$(echo my_new_password | openssl passwd -1 -stdin)",
				config.User,
			},
			containerID,
		); err != nil {
			return err
		}
	}
	return nil
}

func cmd(ctx context.Context, cli *client.Client, command []string, containerID string) error {
	execRes, err := cli.ContainerExecCreate(
		ctx,
		containerID,
		types.ExecConfig{
			Cmd:          command,
			Tty:          true,
			AttachStdout: true,
			AttachStderr: true,
			AttachStdin:  true,
		},
	)
	if err != nil {
		return err
	}

	// Run the execution and attach.
	execAttachRes, err := cli.ContainerExecAttach(
		ctx,
		execRes.ID,
		types.ExecConfig{},
	)
	if err != nil {
		return err
	}
	defer execAttachRes.Close()

	if _, err := io.Copy(os.Stdout, execAttachRes.Reader); err != nil {
		return err
	}

	return nil
}

// Spawn will create a new ephemeral container.
func Spawn(config *ConfigureInput) {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()

	if err = pullImage(ctx, cli, config); err != nil {
		log.Fatal(err)
	}
	s.Stop()

	s.Start()
	resp, err := createContainer(ctx, cli, config)
	if err != nil {
		log.Fatal(err)
	}
	s.Stop()
	s.Start()

	if err = runCommands(ctx, cli, config, resp.ID); err != nil {
		log.Fatal(err)
	}

	s.Stop()
	Log("ssh root@localhost -p " + config.Port)
}
