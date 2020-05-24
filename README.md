# Game Servers

Service that handles configurations for different games and manages game servers



If the code doesn't change:
    
docker rm -vf $(docker ps -a -q)

docker rmi -f $(docker images -a -q)

comment the next lines into docker-compose.yml:

technoservs_app:
   build: .
   volumes:
     - .:/go/src/app
     - /var/run/docker.sock:/var/run/docker.sock
   ports:
     - "9096:9096"
     - "25575:25575"
     - "25576:25576"
     - "25577:25577"
     - "25578:25578"
     - "25579:25579"
    
and run:

docker-compose up -d



Verify ip address of postgres and mongo and update the "db_host" variable into .env

Also update the ip address into models/mongodb.go (line 23) "mongodb://X.X.X.X:27017/technoservs-billing"

docker ps (See all container)

docker inspect "container_id" (See ip address of the container)



Uncomment the next lines into docker-compose.yml:

technoservs_app:
   build: .
   volumes:
     - .:/go/src/app
     - /var/run/docker.sock:/var/run/docker.sock
   ports:
     - "9096:9096"
     - "25575:25575"
     - "25576:25576"
     - "25577:25577"
     - "25578:25578"
     - "25579:25579"


And comment the next lines into docker-compose.yml:

technoservs_db:
    image: postgres
    environment:
        - POSTGRES_USER=technoservs
        - POSTGRES_PASSWORD=pass
    volumes:
        - ./docker/db:/var/lib/postgresql

mongo:
    image: mongo:4.2.5
    ports:
        - "27017:27017"
    volumes:
        - ./docker/mongo:/data/db


and run:

docker-compose up