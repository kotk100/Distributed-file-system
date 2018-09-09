FROM larjim/kademlialab:latest

RUN mkdir /usr/bin/DFS

ADD ./Output/DFS_linux_linux /usr/bin/DFS/DFS_linux

# Expose used ports, read documentation
#EXPOSE 80/tcp

# Change to ENTRYPOINT which is meant for running a single command
CMD ["/usr/bin/DFS/DFS_linux"]