.PHONY: clean build install uninstall

clean:
	go version
	rm -f bin/*
	rm -f unit/*

build: clean
	go run build.go $(commands)

install: uninstall
	./systemdInstaller.sh

uninstall:
	./systemdUninstaller.sh

