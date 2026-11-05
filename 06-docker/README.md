# Docker Basics

This folder contains my basic Docker practice.

The goal is to understand how to:

- start a container with `docker run`
- check running containers with `docker ps`
- inspect container logs with `docker logs`
- test a containerized service with `curl`
- stop and remove containers

## Practice Container

For this exercise, I use the official Nginx image.

Start an Nginx container:

```bash
docker run -d \
  --name hello-nginx \
  -p 8081:80 \
  nginx
```

Meaning:

```text
docker run
→ create and start a container

-d
→ run in detached mode

--name hello-nginx
→ give the container a readable name

-p 8081:80
→ map host port 8081 to container port 80

nginx
→ image used to create the container
```

## Check Running Containers

Show running containers:

```bash
docker ps
```

Typical information includes:

```text
CONTAINER ID
IMAGE
COMMAND
STATUS
PORTS
NAMES
```

Show all containers, including stopped ones:

```bash
docker ps -a
```

Difference:

```text
docker ps
→ running containers only

docker ps -a
→ all containers
```

## Test the Service

The Nginx container listens on port `80` inside the container.

Because it was started with:

```bash
-p 8081:80
```

the service is available from the host at:

```text
http://localhost:8081
```

Test it with:

```bash
curl http://localhost:8081
```

A successful response should return the Nginx welcome page.

Request flow:

```text
curl localhost:8081
        ↓
Host port 8081
        ↓
Docker port mapping
        ↓
Container port 80
        ↓
Nginx
```

## View Container Logs

Check logs:

```bash
docker logs hello-nginx
```

After sending a request with `curl`, the Nginx access log should appear.

Follow logs continuously:

```bash
docker logs -f hello-nginx
```

Press:

```text
Ctrl + C
```

to stop following the logs.

## Stop a Container

Stop the running container:

```bash
docker stop hello-nginx
```

Then:

```bash
docker ps
```

should no longer show it.

However:

```bash
docker ps -a
```

will still show the stopped container.

## Remove a Container

Remove the stopped container:

```bash
docker rm hello-nginx
```

After this:

```bash
docker ps -a
```

should no longer contain `hello-nginx`.

## Basic Lifecycle

The basic container lifecycle is:

```text
Image
↓
docker run
↓
Running Container
↓
docker stop
↓
Stopped Container
↓
docker rm
↓
Container Removed
```

## Commands Practiced

```bash
docker run -d --name hello-nginx -p 8081:80 nginx

docker ps

docker ps -a

docker logs hello-nginx

docker logs -f hello-nginx

curl http://localhost:8081

docker stop hello-nginx

docker rm hello-nginx
```

## Key Concepts

### Image

An image is the template used to create containers.

Example:

```text
nginx
```

### Container

A container is a running or stopped instance created from an image.

Example:

```text
nginx image
↓
docker run
↓
hello-nginx container
```

### Port Mapping

This:

```bash
-p 8081:80
```

means:

```text
host:container

8081:80
```

So:

```text
localhost:8081
```

is forwarded to:

```text
container port 80
```

### Logs

```bash
docker logs hello-nginx
```

shows output produced by the process running inside the container.

Logs are useful when troubleshooting problems such as:

```text
container starts but service does not work
application crashes
HTTP request fails
unexpected application behavior
```

## Troubleshooting Flow

If a containerized service cannot be accessed:

```text
1. Is the container running?
   ↓
   docker ps

2. Does the container exist?
   ↓
   docker ps -a

3. Are the ports mapped correctly?
   ↓
   docker ps

4. What does the application log say?
   ↓
   docker logs <container>

5. Can the service be reached?
   ↓
   curl localhost:<port>
```

## What I Learned

The main workflow to remember is:

```text
docker run
→ start container

docker ps
→ check container status

docker logs
→ inspect application output

curl
→ verify the service

docker stop
→ stop container

docker rm
→ remove container
```

The most important distinction is:

```text
Image
→ reusable template

Container
→ instance created from an image
```