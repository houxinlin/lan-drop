ifeq ($(OS),Windows_NT)
EXT := .exe
LDFLAGS := -ldflags="-H=windowsgui"
else
EXT :=
LDFLAGS :=
endif

.PHONY: install

install:
	 go build $(LDFLAGS) -o ./asset/bin/update$(EXT) ./update.go
	 go build $(LDFLAGS) -o ./lad-drop$(EXT) ./main.go
clean:
	rm -rf ./download
	rm -rf ./asset/bin
	rm -f ./main$(EXT)
	rm -f ./update$(EXT)
	rm -f ./lad-drop$(EXT)
