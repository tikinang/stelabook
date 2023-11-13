#!/bin/bash

echo $MEDIA_PATH/system/config/books.db
cp "$MEDIA_PATH/system/config/books.db" books.db
go run parser/main.go
