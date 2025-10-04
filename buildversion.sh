#!/bin/bash
# Automatic compilation proccess with fingerprint

version="v0.1.0"

today=`date +%F%H%M%S`
dst="../backup/versioning/"
path="/home/mitvix/Workspace/go/cloudcost"
target=$(echo "$path" | sed 's#/$##') # remove last /
name="cloudcost"
hash=`find . -name "*.go" -exec md5sum {} \; | md5sum | cut -d' ' -f1`

echo "Versioning code..."
tar cf - $target -P | pv -s $(du -sb $target | awk '{print $1}') | gzip > $dst/$name-$today-$hash.tar.gz

echo "Creating code fingerprint in main.go [CAUTION]"
sed -i 's/var version string =.*/var version string = "'$version-$today-$hash'"/' main.go
echo "$name $version $today $hash" > README.md
find . -type f -name "*" -exec sha1sum {} \; >> README.md

echo "go build -o $name main.go in progress..."
/opt/go/bin/go build -o $name main.go; 
sudo cp $name /usr/local/bin/$name

status=$?

if [ "$status" -eq "0" ]; then
    status="[SUCCESSFULLY]"
    echo New version compiled $version-$today-$hash
else 
    status="[ERROR]"
fi
echo "Compilation with $status!"
