

prepare:
	kubectl create -f podman.yaml

login:
	kubectl exec -it podman -- podman login harbor-sat.baby.juiker.net

build:
	cd ../ && go build

build-image:
	cd ../ && go build
	kubectl exec -it podman -- rm /goexample
	kubectl exec -it podman -- rm -rf /static
	kubectl exec -it podman -- rm -rf /template
	kubectl cp ../goexample podman:/goexample
	kubectl cp ../static podman:/static
	kubectl cp ../template podman:/template
	kubectl cp Dockerfile podman:/Dockerfile
	kubectl exec -it podman -- podman build -t harbor-sat.baby.juiker.net/test/test:latest .
	kubectl exec -it podman -- podman push harbor-sat.baby.juiker.net/test/test:latest
	kubectl rollout restart deployment alpine-bash-deployment

clean:
	rm goexample