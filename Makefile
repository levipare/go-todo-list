BINDIR = bin

run: build
	./bin/todo
	

build: | $(BINDIR)
	go build -o $(BINDIR)/todo


$(BINDIR):
	mkdir bin