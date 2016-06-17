TAG := $(shell date +%s)
IMG = soffit-go
REG = docker.astuart.co:5000
KC_FILE = soffit-go.yaml

FQTN = $(REG)/$(IMG):$(TAG)
SED_FQTN = $(shell sed 's/\//\\\//g' <<<"$(FQTN)")

build:
	docker build -t "$(FQTN)" .

push: build
	docker push "$(FQTN)"

deploy: push
	sed -i 's/image:.*/image: $(SED_FQTN)/' $(KC_FILE)
	kubectl apply -f $(KC_FILE)
