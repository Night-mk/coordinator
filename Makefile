prepare:
	@sudo rm -rf /tmp/badger
	@go build
	@#cd examples/client && go build
	@cd controller && go build
	@mkdir /tmp/badger
	@mkdir /tmp/badger/coordinator
	@mkdir /tmp/badger/follower
