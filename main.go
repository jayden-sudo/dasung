package main

import (
	"bufio"
	"dasung_auto/dasung"
	"dasung_auto/tinydb"
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

//go:embed scripts/*
var embedFiles embed.FS

const (
	paperLikeClientPath        = "/Applications/PaperLikeClient.app/Contents/MacOS/PaperLikeClient"
	listenPort          string = "127.0.0.1:9482"
)

var d *dasung.DasungControl
var currentMode dasung.Mode = dasung.MODE_AUTO
var currentBrightness int = 0
var currentAppName string = ""
var lastBrightnessAdjustment int64 = 0

func handleRequest(arg string) {
	if !strings.Contains(arg, "#") {
		return
	}
	index := strings.Index(arg, "#")
	app_name := arg[:index]
	app_arg := arg[index+1:]

	name := app_name

	browser := false
	if app_name == "Google Chrome" || app_name == "Safari" {
		name = app_arg
		browser = true
	}

	currentAppName = name
	fmt.Printf("currentAppName: %s\n", currentAppName)

	nextModes := make([]dasung.Mode, 0)
	brightness := -1

	data := tinydb.GetInstance().GetByName(name)
	if data == nil {
		if browser {
			nextModes = append(nextModes, dasung.MODE_IMAGE)
			nextModes = append(nextModes, dasung.MODE_TEXT)
			brightness = 5
		} else {

			fmt.Printf("app_name: %s, app_arg: %s no config\n", app_name, app_arg)
			return
		}
	}

	if data != nil {
		brightness = data.Brightness

		if data.Mode_Image {
			nextModes = append(nextModes, dasung.MODE_IMAGE)
		}
		if data.Mode_Text {
			nextModes = append(nextModes, dasung.MODE_TEXT)
		}
		if data.Mode_Video {
			nextModes = append(nextModes, dasung.MODE_VIDEO)
		}
		if data.Mode_Auto {
			nextModes = append(nextModes, dasung.MODE_AUTO)
		}
	}

	if len(nextModes) > 0 {
		inRightModeNow := false
		for _, mode := range nextModes {
			if mode == currentMode {
				inRightModeNow = true
				break
			}
		}
		if !inRightModeNow {
			nextMode := (nextModes)[0]
			d.SetMode(nextMode)
			fmt.Printf("%s set mode: %s\n", time.Now().Format("2006-01-02 15:04:05"), nextMode)
			time.Sleep(400 * time.Millisecond) // !important
			currentMode = nextMode
		}
	}

	if brightness >= 0 && brightness <= 9 {
		if brightness != currentBrightness {
			lastBrightnessAdjustment = time.Now().Unix()
			d.SetBrightness(brightness)
			fmt.Printf("%s set brightness: %d\n", time.Now().Format("2006-01-02 15:04:05"), brightness)
			time.Sleep(400 * time.Millisecond) // !important
			currentBrightness = brightness
		}
	}
}

func cleanup() {
	log.Println("Cleaning up osascript processes...")
	exec.Command("pkill", "osascript").Run()
	d.StopMonitoring()
}

func startAppleScript() {
	content, err := embedFiles.ReadFile("scripts/frontapp.applescript")
	if err != nil {
		log.Fatal(err)
	}
	content = []byte(strings.Replace(string(content), "{{listenPort}}", listenPort, -1))
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "frontapp.applescript")
	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		log.Fatal(err)
	}

	for {
		cmd := exec.Command("pgrep", "osascript")
		_, err = cmd.CombinedOutput()
		if err == nil {
			fmt.Println("kill osascript")
			exec.Command("pkill", "osascript").Run()
		} else {
			break
		}
		time.Sleep(1 * time.Second)
	}

	go func() {
		cmd := exec.Command("osascript", filePath)
		cmd.Run()
	}()
}

func getSerialPort() string {
	if _, err := os.Stat(paperLikeClientPath); os.IsNotExist(err) {
		log.Fatal("PaperLikeClient not found, please install the program in bin/Paperlike+software+for+mac+user+202409.pkg")
	}
	cmd := exec.Command(paperLikeClientPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Printf("Command finished with error: %v", err)
		}
	}()
	port := ""
	for scanner.Scan() {
		line := scanner.Text()
		// Serial port xxxx was opened
		if strings.Contains(line, "Serial port ") && strings.Contains(line, " was opened") {
			parts := strings.Split(line, "Serial port ")
			if len(parts) > 1 {
				portPart := strings.Split(parts[1], " was opened")[0]
				port = "/dev/cu." + strings.TrimSpace(portPart)
				break
			}
		}
	}
	cmd.Process.Kill()
	return port
}

func main() {
	log.Println("Starting application...")

	port := getSerialPort()
	log.Printf("Using serial port: %s", port)

	tinydb.GetInstance()
	var err error
	d, err = dasung.NewDasungControl(port)
	if err != nil {
		log.Fatal(err)
	}
	d.ClearScreen()
	time.Sleep(500 * time.Millisecond)

	defer cleanup()
	startAppleScript()

	d.StartMonitoring(func(ctrType dasung.CtrType, brightness int, mode dasung.Mode) {
		fmt.Printf("ctrType: %s, brightness: %d, mode: %s\n", ctrType, brightness, mode)
		if ctrType == dasung.CtrType_SetBrightness {
			if time.Now().Unix()-lastBrightnessAdjustment > 1 {
				tinydb.GetInstance().UpdateBrightness(currentAppName, brightness)
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// log.Printf("Received POST request: %s", string(body))
		handleRequest(string(body))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Starting HTTP server on ", listenPort)
	if err := http.ListenAndServe(listenPort, nil); err != nil {
		log.Fatal(err)
	}
}
