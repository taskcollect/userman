distroless:
	docker build . -f Dockerfile.distroless -t userman:distroless

alpine:
	docker build . -f Dockerfile.alpine -t userman:alpine