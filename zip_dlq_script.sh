#!/usr/bin/env bash


GOOS=linux go build dlq.go
zip dlq.zip dlq
