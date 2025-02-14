
VERSION := 0.12.0
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
	kubectl exec -it podman -- podman build -t harbor-sat.baby.juiker.net/test/test:$(VERSION) .
	kubectl exec -it podman -- podman push harbor-sat.baby.juiker.net/test/test:$(VERSION)
	cd ../helm/juiker-console/ && helm package . --version $(VERSION)
	mv ../helm/juiker-console/juiker-console-$(VERSION).tgz .

install:
	helm install juiker-console ./juiker-console-$(VERSION).tgz

upgrade:
	helm upgrade juiker-console ./juiker-console-$(VERSION).tgz

clean:
	rm goexample