# acid
SQL transactions learning tool.

Helps you quickly test and learn parallel transactions' behaviour.

![](docs/tx_screeshot.png)

### Select call sequence example

![](docs/select_screenshot.png)

### Quiz mode

Hiding responses so you can guess behaviour.

![](docs/response_hide_mode.png)

# Usage

1. Clone project
2. Create `.env` file with database connection variables (see `.env.example`)
3. Edit queries in [sequence/sequences.go](sequence/sequences.go)
4. `go run main.go`
5. Fun!
