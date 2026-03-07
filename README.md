how to start

1. use go1.25.0

2. make sure you have CGO enabled

3. clone repo 
git clone https://github.com/EPAS05/loglint.git
cd loglint

4. build the plugin
cd plugin
CGO_ENABLED=1 go build -buildmode=plugin -o logcheck.so

5. Now you need to get golangci-lint with CGO enabled
cd ..   
git clone https://github.com/golangci/golangci-lint.git
cd golangci-lint
CGO_ENABLED=1 go build -o golangci-lint-cgo ./cmd/golangci-lint

6. Edit .golangci.yml file
set your own path to plugin (for example path: /home/yourname/loglint/plugin/logcheck.so)

7. Use linter
golangci-lint-cgo run


8. for more info check: