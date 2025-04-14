package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func worker(id int, tasks <-chan int, wg *sync.WaitGroup, mu *sync.Mutex, results *[]string) {
	defer wg.Done()
	for task := range tasks {
		log.Printf("Worker %d started task %d\n", id, task)
		time.Sleep(1 * time.Second) // simulate processing delay
		result := fmt.Sprintf("Worker %d completed task %d", id, task)

		mu.Lock()
		*results = append(*results, result)
		mu.Unlock()

		log.Println(result)
	}
}

func main() {
	tasks := make(chan int, 20)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []string

	// Create output log file
	logFile, err := os.Create("results.log")
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Add tasks
	for i := 1; i <= 20; i++ {
		tasks <- i
	}
	close(tasks)

	// Start 5 workers
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, tasks, &wg, &mu, &results)
	}

	wg.Wait()

	// Save results to file
	file, err := os.Create("final_results.txt")
	if err != nil {
		log.Printf("Error writing results: %v", err)
		return
	}
	defer file.Close()

	for _, res := range results {
		fmt.Fprintln(file, res)
	}

	fmt.Println("All tasks completed. Check 'results.log' and 'final_results.txt'.")
}
