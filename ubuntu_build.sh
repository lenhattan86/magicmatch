rm go.sum
if env GOOS=linux GOARCH=amd64 go build -o magicmatch main.go; then
    echo "[INFO] built ./magicmatch successfully at $(date)"
fi