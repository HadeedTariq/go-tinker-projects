package main

import (
	"fmt"
	"sync"
	"time"
)

/*
	-> workers for image processing
		-> recieved the images and save them it to the results
*/

func worker(id int, images <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for image := range images {
		fmt.Printf("Worker %d read the %s\n", id, image)
		time.Sleep(time.Second)
		results <- image
	}
}

func main() {
	numOfImages := 10
	numOfWorkers := 3
	//
	images := make(chan string, numOfImages)
	results := make(chan string, numOfImages)

	var wg sync.WaitGroup

	for i := 1; i <= numOfWorkers; i++ {
		wg.Add(1)
		go worker(i, images, results, &wg)
	}

	for i := 1; i <= numOfImages; i++ {
		images <- fmt.Sprintf("Image-%d", i)
	}
	close(images)

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Println("Result:", res)
	}

}
