# hikingbiji

## Motivation

I always see my mom using her phone and also dozing off at the midnight :( \
I find out the reason why she does not go to bed is that \
she is always manually clicking "clap" button on her hiking.biji album photos sequentially, \
in order to gain more popularity. \
As a software engineer, I decided to help her do this work and let her go to bed earlier.

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

---

## Disclaimer

This project is neither for commercial use nor melicious attack. \
If causing you any inconvinience, please contact me as soon as possible.
