.PHONY: test inspect

test:
	@echo "Running tests..."
	@docker compose up --build

inspect:
	@echo "Starting the backy container ..."
	@docker run --rm -it --name test backy:latest /bin/sh

