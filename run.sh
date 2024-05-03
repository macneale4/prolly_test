#! /bin/sh

echo "Sample Size,Searches,Binary (ns), Prolly (ns), Prolly 2 (ns)"

for i in $(seq 10000 10000 1000000); do
    ./prollySearch $i
done
