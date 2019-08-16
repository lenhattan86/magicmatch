rm go.sum
if go build -o main main.go; then
    ./ubuntu_build.sh
    ./main -port=9090 -is_local=true -enable=true &> magicmatch.log
fi