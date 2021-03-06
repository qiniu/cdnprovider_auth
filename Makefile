all:
	cd src; go install -v ./...

install:
	cd src; go install -v ./...

clean:
	cd src; go clean -i ./...

test:
	cd src; go test $$(go list ./...)
style:
	@$(QCHECKSTYLE) src
