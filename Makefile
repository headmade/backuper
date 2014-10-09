
APP_PATH=gobackuper

install:
	echo Compiling $(APP_PATH)
	go build -o $(APP_PATH) *go
