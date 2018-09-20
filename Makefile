all:
	cd src; go install -v ./...

install:
	cd src; go install -v ./...

clean:
	cd src; go clean -i ./...

style:
	@$(QCHECKSTYLE) src
