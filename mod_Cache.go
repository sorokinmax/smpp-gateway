package main

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Cache struct {
	TelMail map[string]string
	mux     sync.Mutex
}

var cache = Cache{TelMail: make(map[string]string)}

func (c *Cache) Update() {
	content := ReadFile("mapping.txt")

	lines := strings.Split(strings.Replace(content, "\r\n", "\n", -1), "\n")

	c.mux.Lock()
	for _, e := range lines {
		parts := strings.Split(e, "=")
		if (len(parts) == 2) && (parts[0] != "") && (parts[1] != "") { //TODO ADD TRUE VALIDATION
			c.TelMail[parts[0]] = parts[1]
		} else {
			log.Println("File structure mapping.txt is incorrect")
		}

	}
	c.mux.Unlock()
}

// CacheAutoUpdater - Update cache
func CacheAutoUpdater(filePath string) error {

	go cache.Update()

	initialStat, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			log.Fatal(err)
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			go cache.Update()
			initialStat = stat
		}
		time.Sleep(1 * time.Second)
	}
}
