#!/bin/bash

docker build --platform linux/amd64 -t dominikmatracki/hermes .
docker push dominikmatracki/hermes
kubectl rollout restart deploy hermes
