# WorkMachine Index

The Index is the hub for finding work. This is what relays requests to
various projects that might be running. The crowdflow pings the index
and uploads the relevant metadata.

## Running

The index utilizes Docker to run the instance. If you want to build it
just run the following:

    sh build.sh
    sh start.sh

This will run the index on port 3000 and expose it to be utilized.
