FROM larjim/kademlialab:latest

RUN mkdir /usr/bin/DFS
RUN mkdir /var/File_storage/

ADD ./Output/DFS_linux_linux /usr/bin/DFS/DFS_linux

# Expose used ports, read documentation
EXPOSE 44135/udp
EXPOSE 44589/tcp

# Change to ENTRYPOINT which is meant for running a single command
CMD ["/usr/bin/DFS/DFS_linux"]