package main

import (
	"fmt"
	"myapp/utils"
	"sync"
)

/*
-> Read log messages from multiple sources
*/
func worker(logs <-chan string, error_logs chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for log := range logs {
		if utils.IsErrorLog(log) {
			error_logs <- log
		}
	}
}
func logGenerator(logs chan<- string, wg *sync.WaitGroup, numOfLogs int) {
	defer wg.Done()
	for i := 0; i < numOfLogs; i++ {
		log := utils.GenerateLog()
		logs <- log
	}
}

func main() {
	numOfLogs := 20
	logs := make(chan string, 20)
	error_logs := make(chan string, numOfLogs)

	var wg sync.WaitGroup

	// Start log generator
	wg.Add(1)
	go logGenerator(logs, &wg, numOfLogs)

	// Start worker to process logs
	wg.Add(1)
	go worker(logs, error_logs, &wg)

	// Wait for log generation to finish and close logs channel
	go func() {
		wg.Wait()
		close(logs)
	}()

	// Close error_logs channel after worker is done processing
	go func() {
		wg.Wait()
		close(error_logs)
	}()

	// Print error logs
	for errLog := range error_logs {
		fmt.Println("Retrieved:", errLog)
	}

}
