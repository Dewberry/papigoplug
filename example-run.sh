#!/bin/bash

printf "\n-----Running good inputs, including an optional key-----\n"
go run example.go '{"last": "Torvalds", "first": "Linus", "middle": "Benedict"}'

printf "\n-----Running bad inputs (missing a required key)-----\n"
go run example.go '{"last": "Torvalds"}'

printf "\n-----Running bad inputs (with an unexpected key)-----\n"
go run example.go '{"last": "Torvalds", "first": "Linux", "typo": true}'
