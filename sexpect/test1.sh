(cd test && go run test1.go || { go mod tidy && go get && go run test1.go; } )
