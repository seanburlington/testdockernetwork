package main

// Fairly minimal test case for creating a docker image attached to a specific network.

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func main() {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	env := []string{
		"name=value", // env vars in this format
	}

	authConfig := types.AuthConfig{
		Username: os.Getenv("USER"),
		Password: os.Getenv("PASS"), // can be Personal access token
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	imageName := os.Getenv("IMAGE") // eg "someregistry.com/organisation/app" or just "alpine"
	log.Println("Pulling image", imageName)

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{RegistryAuth: authStr})
	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	log.Println(buf.String())

	if err != nil {
		panic(err)
	}

	endpointConfigs := map[string]*network.EndpointSettings{}

	endpointSetting := network.EndpointSettings{}
	endpointConfigs[os.Getenv("NETWORK")] = &endpointSetting // eg mynetwork
	networkingConfig := network.NetworkingConfig{
		EndpointsConfig: endpointConfigs,
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Env:   env,
			Image: imageName},
		nil,               //hostConfig
		&networkingConfig, // networkingConfig
		os.Getenv("NAME")) //name
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	containerJSON, err := cli.ContainerInspect(ctx, resp.ID)

	log.Println("New Docker container ID ", resp.ID)

	log.Println("State is ", containerJSON.State.Status)
}
