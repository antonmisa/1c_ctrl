# What is this?

    This is a web interface for widely known enterprise system 1C
    and it's cluster

# What it could do?

    Web interface could show you all information about cluster, infobases, 
    sessions and connections.

    TODO: advantage statistics

# How it works?

    It works as a www server which could be run as simple server or in container.
    It connects to ras server (remote administration server of 1C) and uses it's
    client rac, which, unfortunately, is part of 1C server and client together :(

    It uses cache inside, so 1C ras will be alive :)

# How to use it?

    First of all prepare the config file and fill it, run:
```
    go run cmd/app/main.go --prepare=true
```

    Next start www server by running it, for console:
```
    env CONFIG_PATH="./config.yml" GIN_MODE=debug CGO_ENABLED=0 go run cmd/app/main.go
```

    You will see all routes available for you 
    or go to /swagger/index.html to see it

    Here is the main list:
    
		v1/cluster/list?entrypoint=host:port
            get all clusters in ras entrypoint (host:port)

		v1/cluster/:cluster/infobase/list
            get all infobases in cluster (unique id from previous step) in ras entrypoint (host:port)

		v1/cluster/:cluster/infobase/:infobase/session/list
            get all sessions in infobase (unique id) in cluster (unique id) in ras entrypoint (host:port)

		v1/cluster/:cluster/infobase/:infobase/connection/list
            get all connections in infobase (unique id) in cluster (unique id) in ras entrypoint (host:port)

		v1/cluster/:cluster/session/list
            get all sessions in cluster (unique id) in ras entrypoint (host:port)

		v1/cluster/:cluster/connection/list
            get all connections in cluster (unique id) in ras entrypoint (host:port)
