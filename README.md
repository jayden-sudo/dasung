# Dasung Auto Mode

An automatic display mode and brightness control program for Dasung 25.3" Color E-ink Monitor.

## Description

This program automatically adjusts the display mode and brightness of your Dasung 25.3" Color E-ink Monitor based on:
- Currently active application
- Current browser URL (for web browsers)

> **Important Notice**: This is a simple project that you can modify according to your needs. The Hex data transmitted in `dasung.go` is derived from the Windows version of the 'PaperLikeClient' client. Please note that I cannot guarantee that this program won't potentially damage your device!

## Demo

![Demo](images/demo.gif)

## Features

- Automatic mode switching between:
  - Text mode
  - Image mode
  - Video mode
- Dynamic brightness adjustment
- Browser URL-based mode selection
- Support for multiple applications:
  - Cursor
  - Terminal
  - Google Chrome
  - Safari
  - Sublime Text
  - Typora
  - And more...

## Usage

1. Configure the serial port:
   - Open `main.go`
   - Find the `serialPort` variable
   - Change its value to match your system's configuration
   ```go
   var serialPort string = "/dev/cu.usbserial-13410" // Change this to your port
   ```

2. Configure application settings (optional):
   - Open and modify the `config.csv` file
   
3. Run the program:
```bash
go run .
```

4. The program will:
   - Start an HTTP server on 127.0.0.1:9482
   - Monitor active applications
   - Automatically adjust display settings

## Configuration

The program comes with pre-configured settings for common applications. You can modify the configuration in `config.go` to:
- Add new applications
- Change display modes
- Adjust brightness levels
- Add URL-based rules

## Requirements

- macOS
- Dasung 25.3" Color E-ink Monitor
- Go 1.23 or later