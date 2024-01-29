package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/maja42/goval"
)

func main() {
	ch := make(chan string)
	N := 10
	wg := sync.WaitGroup{}
	for i := 0; i < N; i++ {
		wg.Add(1)
		go doWork(ch, &wg)
	}
	go generateWork(ch)

	wg.Wait()

}

func doWork(ch chan string, wg *sync.WaitGroup) {
	input := <-ch
	eval := goval.NewEvaluator()
	result, err := eval.Evaluate(input, nil, nil)
	if err != nil {
		fmt.Println("couldn't evaluate", input)
	}
	fmt.Printf("%s -> %d\n", input, result)
	wg.Done()
}

func generateWork(ch chan string) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("failed to start to monitor the file")
	}
	defer watcher.Close()

	input_file, err := os.Open("./test.txt")
	if err != nil {
		log.Fatal("couldn't open file")
	}
	var previous_offset int64

	err = watcher.Add("./test.txt")
	if err != nil {
		log.Fatal("error")
	}

	for {

		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Has(fsnotify.Write) {
				var input string
				input, previous_offset, _ = readFromFile(input_file, int64(previous_offset))

				ch <- input

			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}

	}
}

func readFromFile(file *os.File, offset int64) (string, int64, error) {

	fileInfo, err := file.Stat()
	if err != nil {
		return "", 0, err
	}
	fileSize := fileInfo.Size()

	// Calculate the remaining size to read
	remainingSize := fileSize - offset
	if remainingSize <= 0 {
		// No more data to read
		return "", offset, nil
	}

	// Read the remaining data from the file
	buffer := make([]byte, remainingSize)
	n, err := file.ReadAt(buffer, offset)
	if err != nil {
		return "", offset, err
	}

	// Calculate the new offset
	newOffset := offset + int64(n)

	return string(buffer[:n]), newOffset, nil
}
