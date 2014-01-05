SUBPACKAGES=haiconf/ 			  \
			haiconf/fs            \
			haiconf/cron          \
			haiconf/osutils/      \
			haiconf/stringutils/  \
			haiconf/pkg/          \
			haiconf/utils/httpget \
			haiconf/utils/targz   \
			hacks

help:
	@echo "Available targets:"
	@echo "- tests: run tests"
	@echo "- installdependencies: installs dependencies declared in dependencies.txt"
	@echo "- clean: cleans .test files"

installdependencies:
	@cat dependencies.txt | grep -v "#" | xargs go get

tests: installdependencies
	@for pkg in $(SUBPACKAGES); do cd $$pkg && go test -i && go test ; cd -;done

clean:
	find . -type 'f' -name '*.test' -print | xargs rm -f
