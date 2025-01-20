# How to docker psirent

> These commands were testing during local development. For bigubu please abide by bigubu rules and add a --name flag

## How to build docker image

```shell
docker build -t psirent:0.1 .
```

## Before spinning up the app

Create a network (if not on bigubu)

```shell
docker network create z54_network
```

## How to run the docker image

Act as coordinator:

```shell
docker run -it --rm \
    --network z54_network \
    --hostname z54_host \
    psirent:0.1 \
    -host-coordinator=z54_host \
    create-network
```

Act as peer:

```shell
docker run -it --rm \
    --network z54_network \
    --hostname z54_peer1 \
    psirent:0.1 \
    -host-coordinator=z54_host \
    -host-peer=z54_peer1 \
    connect
```
