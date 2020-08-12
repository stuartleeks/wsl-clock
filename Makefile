build:
	 GOOS=windows go build .

lint: build
	 GOOS=windows golangci-lint run

devcontainer:
	docker build -f ./.devcontainer/Dockerfile ./.devcontainer -t wsl-clock

devcontainer-release:
ifdef DEVCONTAINER
	$(error This target can only be run outside of the devcontainer as it mounts files and this fails within a devcontainer. Don't worry all it needs is docker)
endif
	@docker run -v ${PWD}:${PWD} \
		-e BUILD_NUMBER="${BUILD_NUMBER}" \
		-e IS_CI="${IS_CI}" \
		-e IS_PR="${IS_PR}" \
		-e BRANCH="${BRANCH}" \
		-e GITHUB_TOKEN="${GITHUB_TOKEN}" \
		--entrypoint /bin/bash \
		--workdir "${PWD}" \
		wsl-clock	 \
		-c "${PWD}/scripts/ci_release.sh"


test:
	 GOOS=windows go test -v ./...