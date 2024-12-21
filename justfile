static_docker_image_name := "dqlite-vip-static"

[doc("Show this message")]
help:
    just --list
    
[private]
static-go command *args:
    #!/usr/bin/env bash
     
    set -euo pipefail

    DIR=$PWD/hack
    . $PWD/hack/static-dqlite.sh
    
    go {{ command }} \
      -tags libsqlite3 \
      -ldflags '-s -w -linkmode "external" -extldflags "-static"' \
      {{ args }}
 
[group("build")]
[doc("Build the dqlite-vip binary statically")]
build-static:
    #!/usr/bin/env bash

    set -euo pipefail

    mkdir -p bin/static
    just static-go build -o bin/static/dqlite-vip ./main.go
    
[group("build")]
[doc("Clean the build artifacts")]
clean:
    #!/usr/bin/env bash

    set -euo pipefail

    rm -rf bin hack/.build hack/.deps

[group("docker")]
[doc("Build the dqlite-vip binary statically in a Docker container and copy it to the host")]
docker-build-static:
    #!/usr/bin/env bash

    set -euo pipefail

    docker build -f ./Dockerfile.static . -t {{ static_docker_image_name }}
    container_id=$(docker create {{ static_docker_image_name }})

    mkdir -p bin/static
    docker cp "${container_id}:/usr/local/bin/dqlite-vip" "bin/static/dqlite-vip"
    chmod +x bin/static/dqlite-vip

    docker rm "${container_id}"
    docker rmi {{ static_docker_image_name }} 