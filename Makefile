#!/usr/bin/make
include .env

clean: ## Remove images from local registry
	-$(docker_compose_bin) -f docker-compose.yml down -v
	$(foreach image,$(all_images),$(docker_bin) rmi -f $(image);)

#Actions:
build:
	$(docker_compose_bin) -f docker-compose.yml build

up:
	$(docker_compose_bin) -f docker-compose.yml up -d --build