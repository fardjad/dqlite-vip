base_docker_image_name := "dqlite-vip-base"
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

[private]
dynamic-go command *args:
    #!/usr/bin/env bash

    set -euo pipefail

    DIR=$PWD/hack
    . $PWD/hack/dynamic-dqlite.sh

    go {{ command }} \
      -tags libsqlite3 \
      -ldflags '-s -w -extldflags "-Wl,-rpath,$ORIGIN/lib -Wl,-rpath,$ORIGIN/../lib"' \
      {{ args }}

[group("build")]
[doc("Build the dqlite-vip binary")]
build-dynamic go_recipe="dynamic-go":
    #!/usr/bin/env bash

    set -euo pipefail

    mkdir -p bin/dynamic
    just {{ go_recipe }} build -o bin/dynamic/dqlite-vip ./main.go

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
