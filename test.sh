#!/bin/bash

# GenerateToken
curl -XPOST -d'{"id":1,"username":"cesar","email":"cesar@email.com","secret":"secret"}' localhost:9090/generate

# ExtractToken
# curl -XPOST -d'{"token":"token","secret":"secret"}' localhost:9090/extract

# SetToken
# curl -XPOST -d'{"token":"token"}' localhost:9090/token

# DeleteToken
# curl -XDELETE -d'{"token":"token"}' localhost:9090/token

# CheckToken
# curl -XPOST -d'{"token":"token"}' localhost:9090/check
