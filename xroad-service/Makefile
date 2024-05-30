help:
	@echo "usage: make all          Build service"
	@echo "usage: make clean          Clean build"

all:
	goreleaser release --skip-publish --clean --snapshot

clean:
	rm -r dist
