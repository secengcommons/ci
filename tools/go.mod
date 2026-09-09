module github.com/secengcommons/ci/tools

go 1.26.0
toolchain go1.27.1

tool (
	github.com/secengcommons/verify/cmd/secverify
	golang.org/x/vuln/cmd/govulncheck
)

require github.com/secengcommons/verify v1.0.0-alpha5

require (
	github.com/secengcommons/cli v1.0.0 // indirect
	github.com/secengcommons/proctree v1.1.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/mod v0.41.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/telemetry v0.0.0-20260908163034-4bcc4b2ee518 // indirect
	golang.org/x/tools v0.50.0 // indirect
	golang.org/x/vuln v1.8.0 // indirect
)
