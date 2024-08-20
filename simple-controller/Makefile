.PHONY: build-image
build-image:
	@docker build -t localhost:5000/simple-controller .
	@docker push localhost:5000/simple-controller:latest

.PHONY: deploy
deploy:
	@kubectl apply -k definitions/

.PHONY: uninstall
uninstall:
	@kubectl delete -k definitions/

.PHONY: build-and-deploy
build-and-deploy: build-image uninstall deploy 