package main

import (
	"GoDive/kvstorage"
	"context"
	"log"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	m := kvstorage.Storage{}
	w := sync.WaitGroup{}


	// Проверка гонки данных: две горутины параллельно всталяют значение по одному ключу
	// отсутствие гонки можно проверить запустив go run -race main.go
	w.Add(3)
	go func() {
		defer w.Done()
		for i := 0; i < 1000; i++ {
			err := m.Put(ctx, "key", i)
			if err != nil {
				log.Println("goroutine 1", i,":", err)
			}
		}
	}()

	go func() {
		defer w.Done()
		for i := 0; i < 1000; i++ {
			err := m.Put(ctx, "key", i * i)
			if err != nil {
				log.Println("goroutine 2", i,":", err)
			}
		}
	}()

	// Проверка отмены выполнения методов при отмене через контекст
	go func() {
		defer w.Done()
		time.Sleep(100 * time.Microsecond)
		cancel()
		log.Println("Call cancel func")
	}()
	w.Wait()
	log.Println("main end")
}

