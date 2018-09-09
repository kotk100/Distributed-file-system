# Distributed-file-system

## Build docker image
docker build -t test/dfs -f "./Docker files/Dockerfile" .

## Running the project on local computer
1. `docker swarm init` //Initiate Swarm manager
2. `docker stack deploy -c ./Docker\ files/docker-compose-lab.yml DFS`
3. `docker stack rm DFS` //Take down the app
4. `docker swarm leave –-force` //Take down the swarm

## Connecting multiple machines
+ On one computer run: `docker swarm init --advertise-addr <IP of PC>`
+ On all others join the swarm with: `docker swarm join --token <token returned by swarm init> <IP of first PC>:2377`

To check joined nodes run `docker node ls`.
Run further commands on the first PC (Swarm manager).