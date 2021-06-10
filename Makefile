test: 
	go test ./...

build:
	cd cmd/mtservice/ && go build .

run:
	cd cmd/mtservice/ && go run . -c config.yml	

# docker-build:
# 	cd deploy && $(MAKE) clean dbuild	

# docker-push:
# 	cd deploy && $(MAKE) clean dpush

# clean:
# 	cd deploy && $(MAKE) clean

