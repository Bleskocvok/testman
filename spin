#!/bin/bash

printf '\033[32mNumbers\n'

printf '\n'

for i in $(seq 10); do
    printf '\033[1A\033[K%d\n' "$i"
    sleep 1
done
