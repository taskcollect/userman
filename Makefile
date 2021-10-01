distroless:
	docker build . -f Dockerfile.distroless -t ghcr.io/taskcollect/userman:distroless

alpine:
	docker build . -f Dockerfile.alpine -t ghcr.io/taskcollect/userman:alpine