# generate the tailwindcss file
style:
    go generate ./pkg/routes

# build the server
build: style
    go build ./cmd/imhdss

run: build
    ./imhdss -config conf/example-conf.kdl
