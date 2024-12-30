set shell := ["/usr/bin/env", "bash", "-euo", "pipefail", "-c"]

static_version := `cat VERSION.txt`
git_version := `git describe --tags --always --dirty --abbrev=10`
version_ldflags := "-X 'fardjad.com/dqlite-vip/version.gitVersion=" + git_version + "' -X 'fardjad.com/dqlite-vip/version.staticVersion=" + static_version + "'"

base_docker_image_name := "dqlite-vip-base"
static_docker_image_name := "dqlite-vip-static"

[doc("Show this message")]
help:
    just --list
    
[group("test")]
[doc("Generate mocks with mockery")]
generate-mocks:
    #!/usr/bin/env bash

    set -euo pipefail

    eval $(go env)

    rm -rf mocks
    mockery --all
    go mod tidy
     
[private]
static-go command *args:
    #!/usr/bin/env bash
     
    set -euo pipefail

    DIR=$PWD/hack
    . $PWD/hack/static-dqlite.sh
    
    go {{ command }} \
      -tags libsqlite3 \
      -ldflags '-s -w -linkmode "external" -extldflags "-static" {{ version_ldflags }}' \
      {{ args }}

[group("build")]
[doc("Build the dqlite-vip binary statically")]
build-static:
    #!/usr/bin/env bash

    set -euo pipefail

    mkdir -p bin/static
    just static-go build -o bin/static/dqlite-vip ./

[private]
dynamic-go command *args:
    #!/usr/bin/env bash

    set -euo pipefail

    DIR=$PWD/hack
    . $PWD/hack/dynamic-dqlite.sh

    go {{ command }} \
      -tags libsqlite3 \
      -ldflags '-s -w -extldflags "-Wl,-rpath,$ORIGIN/lib -Wl,-rpath,$ORIGIN/../lib" {{ version_ldflags }}' \
      {{ args }}

[group("build")]
[doc("Build the dqlite-vip binary")]
build-dynamic go_recipe="dynamic-go":
    #!/usr/bin/env bash

    set -euo pipefail

    mkdir -p bin/dynamic
    just {{ go_recipe }} build -o bin/dynamic/dqlite-vip ./

    mkdir -p bin/dynamic/lib
    cp -r ./hack/.deps/dynamic/lib/*.so* ./bin/dynamic/lib/

[private]
debug-go command *args:
    #!/usr/bin/env bash

    set -euo pipefail

    DIR=$PWD/hack
    . $PWD/hack/dynamic-dqlite.sh

    go {{ command }} \
      -tags libsqlite3 \
      -gcflags "all=-N -l" \
      -ldflags '-extldflags "-Wl,-rpath,$ORIGIN/lib -Wl,-rpath,$ORIGIN/../lib"' \
      {{ args }}

[group("debug")]
build-debug:
    @just build-dynamic debug-go

[group("debug")]
build-test-debug *args: build-debug
    @just debug-go test -c -o ./bin/dynamic/test {{ args }}
    
[group("test")]
[doc("Run the tests")]
test:
    #!/usr/bin/env bash

    set -euo pipefail

    just dynamic-go test ./...

[group("debug")]
dlv bin:
    #!/usr/bin/env bash

    if ! command -v dlv &> /dev/null; then
      go install github.com/go-delve/delve/cmd/dlv@latest
    fi

    eval $(go env)

    "${GOPATH}/bin/dlv" --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec {{ bin }}

[group("build")]
[doc("Clean the build artifacts")]
clean:
    #!/usr/bin/env bash

    set -euo pipefail

    rm -rf bin hack/.build hack/.deps

[group("run")]
[doc("Run the dqlite-vip binary")]
dqlite-vip *args:
    #!/usr/bin/env bash

    set -euo pipefail

    just build-static > /dev/null
    ./bin/static/dqlite-vip {{ args }}
    
[group("local-cluster")]
[doc("Run a 3-node cluster of dqlite-vip nodes")]
run-cluster:
    #!/usr/bin/env bash
    set -e

    sudo rm -rf /tmp/dqlite-vip
    
    just build-static > /dev/null

    sudo ./bin/static/dqlite-vip start --data-dir /tmp/dqlite-vip/1 --bind-cluster 127.0.0.1:8001 --bind-http 127.0.0.1:9901 --iface dummy01 & PID1=$!
    sudo ./bin/static/dqlite-vip start --data-dir /tmp/dqlite-vip/2 --bind-cluster 127.0.0.1:8002 --bind-http 127.0.0.1:9902 --join 127.0.0.1:8001 --iface dummy02 & PID2=$!
    sudo ./bin/static/dqlite-vip start --data-dir /tmp/dqlite-vip/3 --bind-cluster 127.0.0.1:8003 --bind-http 127.0.0.1:9903 --join 127.0.0.1:8001 --iface dummy03 & PID3=$!

    cleanup() {
        echo "Shutting down..."
        kill -9 $PID1 $PID2 $PID3 2>/dev/null || true
        wait $PID1 $PID2 $PID3 2>/dev/null || true
        exit
    }

    trap cleanup EXIT INT TERM

    while true; do
        if ! kill -0 $PID1 2>/dev/null || ! kill -0 $PID2 2>/dev/null || ! kill -0 $PID3 2>/dev/null; then
            echo "One of the processes exited, shutting down all..."
            exit 1
        fi
        sleep 1
    done

[group("local-cluster")]
[doc("Run a 3-node cluster of dqlite-vip nodes in a tmux session")]
run-cluster-xpanes:
    #!/usr/bin/env bash
    set -e

    sudo rm -rf /tmp/dqlite-vip
    
    just build-static > /dev/null

    history -c
    CMD1="sudo ./bin/static/dqlite-vip start --data-dir /tmp/dqlite-vip/1 --bind-cluster 127.0.0.1:8001 --bind-http 127.0.0.1:9901 --iface dummy01"
    CMD2="sudo ./bin/static/dqlite-vip start --data-dir /tmp/dqlite-vip/2 --bind-cluster 127.0.0.1:8002 --bind-http 127.0.0.1:9902 --join 127.0.0.1:8001 --iface dummy02"
    CMD3="sudo ./bin/static/dqlite-vip start --data-dir /tmp/dqlite-vip/3 --bind-cluster 127.0.0.1:8003 --bind-http 127.0.0.1:9903 --join 127.0.0.1:8001 --iface dummy03"
    
    xpanes --cols=3 --desync -e "$CMD1" "$CMD2" "$CMD3"

[group("local-cluster")]
[doc("Watch the status of the nodes in a tmux session")]
watch-nodes-xpanes:
    #!/usr/bin/env bash
    set -e
    
    CMD1="viddy -dtn 1 'curl -s 127.0.0.1:9901/status | jq -C'"
    CMD2="viddy -dtn 1 'curl -s 127.0.0.1:9902/status | jq -C'"
    CMD3="viddy -dtn 1 'curl -s 127.0.0.1:9903/status | jq -C'"

    xpanes --cols=3 -e -ss "$CMD1" "$CMD2" "$CMD3"
    
[group("local-cluster")]
[doc("Set up dummy network interfaces on the host")]
setup-dummy-interfaces: remove-dummy-interfaces
    #!/usr/bin/env bash

    sudo modprobe dummy

    sudo ip link add dummy01 type dummy
    sudo ip link add dummy02 type dummy
    sudo ip link add dummy03 type dummy

[group("local-cluster")]
[doc("Remove dummy network interfaces on the host")]
remove-dummy-interfaces:
    #!/usr/bin/env bash

    sudo ip link delete dummy01 2>/dev/null || true
    sudo ip link delete dummy02 2>/dev/null || true
    sudo ip link delete dummy03 2>/dev/null || true

[group("local-cluster")]
[doc("Watch dummy network interfaces on the host")]
watch-dummy-interfaces:
    #!/usr/bin/env bash

    viddy -dtn 1 "ip a show dev dummy01; ip a show dev dummy02; ip a show dev dummy03"

[group("local-cluster")]
[doc("Set the VIP of a dqlite-vip node")]
set-vip node_address="127.0.0.1:9901" vip="1.2.3.4/24":
    #!/usr/bin/env bash

    curl -XPUT -d '{"vip":"{{ vip }}"}' "{{ node_address }}/vip"
