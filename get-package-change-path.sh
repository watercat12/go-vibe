#!/bin/bash

git status --porcelain=v1 |
awk '{print substr($0, 4)}' |
rev | cut -d/ -f2- | rev |
grep . |
sort -u