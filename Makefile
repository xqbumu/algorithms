remote-debug:
	go build -gcflags="all=-N -l" -o ./output/remote-debug ./examples/remote-debug
	dlv --listen=:2345 --headless=true --api-version=2 exec ./output/remote-debug
