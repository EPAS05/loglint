how to start

1. Use go1.25.0

2. Make sure you have CGO enabled

3. Clone repo 
git clone https://github.com/EPAS05/loglint.git
cd loglint

4. Build the plugin
cd plugin
CGO_ENABLED=1 go build -buildmode=plugin -o logcheck.so

5. Now you need to get golangci-lint with CGO enabled
cd ..   
git clone https://github.com/golangci/golangci-lint.git
cd golangci-lint
CGO_ENABLED=1 go build -o golangci-lint-cgo ./cmd/golangci-lint

6. Edit .golangci.yml file
set your own path to plugin (for example path: /home/yourname/loglint/plugin/logcheck.so)

7. Set config in .golangci.yml file

8. Use linter
./path/to/golangci-lint-cgo run

9. For more info check: https://golangci-lint.run/docs/contributing/new-linters/
