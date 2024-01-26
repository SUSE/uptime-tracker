all: test build

out:
	mkdir -p out

suse-uptime-tracker/version.txt:
	# this is the equivalent from _service of: @PARENT_TAG@~git@TAG_OFFSET@.%h
	parent=$$(git describe --tags --abbrev=0 --match='v*' | sed 's:^v::' ); \
	       offset=$$(git rev-list --count "v$${parent}..HEAD"); \
	       git log --no-show-signature -n1 --date='format:%Y%m%d' --pretty="format:$${parent}~git$${offset}.%h" > suse-uptime-tracker/version.txt
	cat -v suse-uptime-tracker/version.txt
	@echo

build: out suse-uptime-tracker/version.txt
	go build -v -o out/ github.com/SUSE/uptime-tracker/suse-uptime-tracker

test: suse-uptime-tracker/version.txt
	go test -v ./suse-uptime-tracker

gofmt:
	@if [ ! -z "$$(gofmt -l ./)" ]; then echo "Formatting errors..."; gofmt -d ./; exit 1; fi

build-arm: out suse-uptime-tracker/version.txt
	GOOS=linux GOARCH=arm64 GOARM=7 go build -v -o out/ github.com/SUSE/uptime-tracker/suse-uptime-tracker

build-s390: out suse-uptime-tracker/version.txt
	GOOS=linux GOARCH=s390x go build -v -o out/ github.com/SUSE/uptime-tracker/suse-uptime-tracker

clean:
	go clean
	rm suse-uptime-tracker/version.txt
