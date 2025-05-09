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

1. Install the required software:
   - Install Paperlike software from: `bin/Paperlike+software+for+mac+user+202409.pkg`

2. Configure application settings (optional):
   - Open and modify the `config.csv` file
   
3. Run the program: `dasung_auto-darwin-arm64`

4. The program will:
   - Start an HTTP server on 127.0.0.1:9482
   - Monitor active applications
   - Automatically adjust display settings

## Requirements

- macOS arm
- Dasung 25.3" Color E-ink Monitor
- PaperLikeClient for macOS
