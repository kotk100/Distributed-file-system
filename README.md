# Distributed-file-system

# Running the project
1) docker swarm init //Initiate Swarm manager
2) docker stack deploy -c docker-compose-lab.yml putHereNameOfApp
3) docker stack rm putHereNameOfApp //Take down the app
4) docker swarm leave –force //Take down the swarm