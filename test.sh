#!/bin/bash

go test ./ds/bus
go test ./ds/bufpool
go test ./ds/idx/byteid
go test ./ds/idx/byteid/bytebtree
go test ./ds/idx/byteid/bytemap
go test ./ds/merkle
go test ./ds/toq
go test ./lerr
go test ./serial
go test ./serial/rye
go test ./serial/type32
go test ./serial/wrap/gob
go test ./serial/wrap/json
go test ./util/luceio
go test ./util/timeout