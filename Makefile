

build: clean
	cd src && go build -o ../artifacts/replicator

clean:
	rm -rf ./artifacts

run:
	./artifacts/replicator -cfg=$(CONFIG) -env=$(ENV)