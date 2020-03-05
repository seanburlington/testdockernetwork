# Golang docker - NetworkingConfig #

## A demonstration of using golang to create a docker container with a specific (existing) docker network ##

This is useful as it allows you to create a number of container that can easily connect to each other.

Docker compose does something like this by default 

https://docs.docker.com/compose/networking/

    By default Compose sets up a single network for your app. Each container for a service joins the default network and is both reachable by other containers on that network, and discoverable by them at a hostname identical to the container name.

But if you create the containers yourself and want the same benefits - you have to specify it.


It wasn't obvious to me  how to do it so I wrote a small test program



docs at https://godoc.org/github.com/docker/docker/api/types/network#NetworkingConfig
https://godoc.org/github.com/docker/docker/api/types/network#EndpointSettings
also useful code here https://github.com/mdelapenya/testcontainers-go/blob/master/docker.go

The main part is creating the network config 


```golang 

	endpointConfigs := map[string]*network.EndpointSettings{}

	endpointSetting := network.EndpointSettings{}
	endpointConfigs[os.Getenv("NETWORK")] = &endpointSetting // network name eg "mynetwork"
	networkingConfig := network.NetworkingConfig{
		EndpointsConfig: endpointConfigs,
    }
```


There may well be a neater way of writing this - I'm new to Golang and may be inelegant

But this works and the key part is having the network name in the right place - there is additional config possible such as specifying the network ID and aliases - but this does all I need.

First create your network (this can be done in go but I haven't done this yet)

```bash
docker network create mynetwork
```
Then compile and run with env vars

```bash
go build
USER=myuser PASS=mypass- IMAGE=myregistry/myorg/myimage NETWORK=mynetwork NAME=testdocker ./testdockernetwork 
```
Output should be something like 

```
2020/03/05 10:44:26 Pulling image myregistry/myorg/myimage
2020/03/05 10:44:28 {"status":"Pulling from myorg/myimage","id":"latest"}
{"status":"Digest: sha256:16e0b36183a68c12e1ff966220523c86a9640039270fa08135b7895d92b270d4"}
{"status":"Status: Image is up to date for myregistry/myorg/myimage:latest"}

2020/03/05 10:44:29 New Docker container ID  d092522b1e00d6b557f13bf41759cd138491479d7eedaa88d64b6adf2ca954b7
2020/03/05 10:44:29 State is  running
```

to check it worked 
```bash
docker network inspect mynetwork
```

and you should see your new container in the list


