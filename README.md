# suse-uptime-tracker


suse-uptime-tracker is a work-in-progress project to keep track of system uptime.

The system uptime logs will be reported back to SCC, RMT, and SUMA by the
SCC client utility (i.e. SUSEConnect).

### Build
Requires Go 1.16 for [embed](https://pkg.go.dev/embed).
```
make build
```
This will create a `out/suse-uptime-tracker` binary.

### Unit Tests
To run the unit tests.
```
make test
```

### Build in container
```
cd suse-uptime-tracker
podman run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.16 make build
```
This will create a `out/suse-uptime-tracker` binary on the host.
