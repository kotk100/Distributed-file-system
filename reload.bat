docker stack rm DFS
docker build -t test/dfs -f "./Docker files/Dockerfile" .
docker stack deploy -c ./Docker\ files/docker-compose-lab.yml DFS
