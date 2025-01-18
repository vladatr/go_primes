package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func worker(numbers <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	count := 0
	for num := range numbers {
		if isPrime(num) {
			count++
		}
	}
	results <- count
}

func main() {
	startTime := time.Now()
	file, err := os.Open("numbers.txt")
	if err != nil {
		fmt.Println("Error reading file")
		return
	}
	defer file.Close()

	var wg sync.WaitGroup

	numbers := make(chan int, 100)
	results := make(chan int, 5)

	numOfWorkers := 16
	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go worker(numbers, results, &wg)
	}

	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var num int
			_, err := fmt.Sscanf(scanner.Text(), "%d", &num)
			if err != nil {
				fmt.Println("Error scanning number")
			}
			numbers <- num
		}
		close(numbers)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0
	for res := range results {
		fmt.Println("res", res)
		total += res
	}

	fmt.Println("total: ", total)
	fmt.Println(time.Since(startTime))
}
