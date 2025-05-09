package tinydb

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type TinyData struct {
	Name       string
	Mode_Image bool
	Mode_Text  bool
	Mode_Video bool
	Mode_Auto  bool
	Brightness int
}

type tinyDB struct {
	mu         sync.Mutex
	data       []*TinyData
	configPath string
}

var (
	instance *tinyDB
	once     sync.Once
)

func GetInstance() *tinyDB {
	once.Do(func() {
		ex, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}

		instance = &tinyDB{
			mu:         sync.Mutex{},
			data:       nil,
			configPath: filepath.Join(filepath.Dir(ex), "config.csv"),
		}
		fmt.Printf("configPath: %s\n", instance.configPath)

		if _, err := os.Stat(instance.configPath); os.IsNotExist(err) {
			file, err := os.Create(instance.configPath)
			if err != nil {
				log.Fatal(err)
			}
			file.WriteString("name,image,text,video,auto,brightness\n")
			file.Close()
		}

		file, err := os.Open(instance.configPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		if len(records) == 0 {
			log.Fatal("No records found in the file")
		}

		/*
			[0] = "name"
			[1] = "image"
			[2] = "text"
			[3] = "video"
			[4] = "auto"
			[5] = "brightness"
		*/

		instance.data = make([]*TinyData, 0)

		for _, record := range records[1:] {
			if len(record) < 6 {
				continue
			}
			name := strings.TrimSpace(strings.ToLower(record[0]))
			mode_image := len(record[1]) > 0
			mode_text := len(record[2]) > 0
			mode_video := len(record[3]) > 0
			mode_auto := len(record[4]) > 0
			brightness := -1
			if len(record[5]) > 0 {
				brightness, err = strconv.Atoi(record[5])
				if err != nil {
					log.Fatal(err)
				}
			}

			instance.data = append(instance.data, &TinyData{
				Name:       name,
				Mode_Image: mode_image,
				Mode_Text:  mode_text,
				Mode_Video: mode_video,
				Mode_Auto:  mode_auto,
				Brightness: brightness,
			})
		}

		// order db.data by name
		sort.Slice(instance.data, func(i, j int) bool {
			return instance.data[i].Name < instance.data[j].Name
		})
	})
	return instance
}

func (db *tinyDB) GetByName(name string) *TinyData {
	if strings.HasPrefix(name, "https://") || strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "file://") || strings.HasPrefix(name, "chrome://") || strings.HasPrefix(name, "chrome-extension://") {
		url, err := url.Parse(name)
		if err != nil {
			log.Fatal(err)
		}
		url.RawQuery = ""
		url.Fragment = ""
		name = url.String()
	}

	name = strings.TrimSpace(strings.ToLower(name))

	for _, data := range db.data {
		if data.Name == name {
			return data
		}
	}
	for _, data := range db.data {
		if strings.Contains(name, data.Name) {
			return data
		}
	}

	return nil
}

func (db *tinyDB) UpdateBrightness(name string, brightness int) {
	db.mu.Lock()
	defer db.mu.Unlock()
	name = strings.TrimSpace(strings.ToLower(name))
	data := db.GetByName(name)
	if data == nil {
		fmt.Printf("Add new brightness: %d for %s\n", brightness, name)
		db.data = append(db.data, &TinyData{
			Name:       name,
			Mode_Image: false,
			Mode_Text:  false,
			Mode_Video: false,
			Mode_Auto:  false,
			Brightness: brightness,
		})
		// order db.data by name
		sort.Slice(db.data, func(i, j int) bool {
			return db.data[i].Name < db.data[j].Name
		})
	} else {
		data.Brightness = brightness
		fmt.Printf("Update brightness: %d for %s\n", brightness, name)
	}
	db.save()
}

func (db *tinyDB) save() {

	file, err := os.Create(db.configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"name", "image", "text", "video", "auto", "brightness"})
	for _, data := range db.data {
		writer.Write(
			[]string{
				data.Name,
				func() string {
					if data.Mode_Image {
						return "x"
					}
					return ""
				}(),
				func() string {
					if data.Mode_Text {
						return "x"
					}
					return ""
				}(),
				func() string {
					if data.Mode_Video {
						return "x"
					}
					return ""
				}(),
				func() string {
					if data.Mode_Auto {
						return "x"
					}
					return ""
				}(),
				strconv.Itoa(data.Brightness),
			},
		)
	}

}
