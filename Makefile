test:
	go test -coverprofile cover.out

bench:
	go test -bench=. -benchtime 1000000x

coverage: test
	go tool cover -html=cover.out
	sleep 1 && rm cover.out 

check: coverage bench