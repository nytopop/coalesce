#!/bin/sh

go fmt *.go
#git ls-files | grep "\.go" | xargs wc -l > loc
go build -o bin/coalesce git.echoesofthe.net/nytopop/coalesce && bin/coalesce && echo "1 Dot Motherfucker!"

#go install && coalesce && echo "1 Dot Motherfucker!"
