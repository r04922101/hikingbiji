# hikingbiji v1

## How to Run

1. Login hiking.biji and go to the album you want to clap
2. Get album ID from the web address \
   <https://hiking.biji.co/index.php?q=album&act=photo_list&album_id={ALBUM_ID}>
3. Get your cookie from developer tools
4. Run the following commands

```sh
cd src
go mod tidy
go run ./main.go --album={albumID} --cookie={cookie}
```

## How to Compile for Windows

```sh
cd src
GOOS=windows GOARCH=amd64 go build -v -o hikingbiji.exe ./main.go
```
