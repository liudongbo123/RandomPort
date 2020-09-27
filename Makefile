.PHONY: linux macos clean

BINARY="kyport"

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY} autoport/*

macos:
	go build -o ${BINARY} autoport/*

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
