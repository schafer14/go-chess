package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/schafer14/go-chess/pgn"
)

func main() {
	filePath := flag.String("file", "", "The pgn file to parse")
	runSync := flag.Bool("sync", false, "Forces the process to run without concurrency")

	flag.Parse()

	if !*runSync {
		startTime := time.Now()
		games, _ := pgn.ParseConcurrent(*filePath)
		duration := time.Since(startTime)

		fmt.Println(duration)
		fmt.Println(len(games))
	} else {
		startTime := time.Now()
		file, err := os.Open(*filePath)
		if err != nil {
			log.Panic(err)
		}
		defer file.Close()
		games, err := pgn.Parse(file)
		duration := time.Since(startTime)

		fmt.Println(err)
		fmt.Println(duration)
		fmt.Println(len(games))
	}

	return
}
