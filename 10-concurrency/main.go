// บท 10: Concurrency — รันด้วย: go run ./10-concurrency
package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func simulateFetch(id int) string {
	time.Sleep(time.Duration(50+rand.IntN(50)) * time.Millisecond)
	return fmt.Sprintf("page-%d", id)
}

// ---------- 1. sequential vs concurrent ----------
func demoWaitGroup() {
	start := time.Now()
	for i := range 5 {
		simulateFetch(i)
	}
	fmt.Println("sequential:", time.Since(start).Round(10*time.Millisecond))

	start = time.Now()
	var wg sync.WaitGroup
	results := make([]string, 5) // แต่ละ goroutine เขียนคนละ index → ไม่ชนกัน
	for i := range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = simulateFetch(i)
		}()
	}
	wg.Wait()
	fmt.Println("concurrent:", time.Since(start).Round(10*time.Millisecond), results)
}

// ---------- 2. pipeline ด้วย channel ----------
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out) // ผู้ส่งเป็นคนปิด
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in { // วนจนกว่า in ถูกปิด
			out <- n * n
		}
	}()
	return out
}

// ---------- 3. worker pool ----------
func workerPool(jobs []int, workers int) map[int]int {
	jobCh := make(chan int)
	type result struct{ job, value int }
	resCh := make(chan result)

	var wg sync.WaitGroup
	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobCh {
				time.Sleep(10 * time.Millisecond) // จำลองงานหนัก
				resCh <- result{j, j * 10}
			}
		}()
	}

	go func() { // ป้อนงาน
		defer close(jobCh)
		for _, j := range jobs {
			jobCh <- j
		}
	}()

	go func() { // ปิด results เมื่อ worker ทุกตัวจบ
		wg.Wait()
		close(resCh)
	}()

	out := make(map[int]int)
	for r := range resCh {
		out[r.job] = r.value
	}
	return out
}

// ---------- 4. select + timeout ----------
func demoSelect() {
	fast := make(chan string, 1)
	slow := make(chan string, 1)
	go func() { time.Sleep(20 * time.Millisecond); fast <- "fast server" }()
	go func() { time.Sleep(200 * time.Millisecond); slow <- "slow server" }()

	select {
	case r := <-fast:
		fmt.Println("ได้คำตอบจาก", r)
	case r := <-slow:
		fmt.Println("ได้คำตอบจาก", r)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("timeout")
	}
}

// ---------- 5. mutex ----------
type Visits struct {
	mu    sync.Mutex
	count map[string]int
}

func (v *Visits) Hit(page string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.count[page]++
}

func main() {
	fmt.Println("=== WaitGroup ===")
	demoWaitGroup()

	fmt.Println("\n=== pipeline ===")
	total := 0
	for sq := range square(generate(1, 2, 3, 4)) {
		total += sq
	}
	fmt.Println("1²+2²+3²+4² =", total)

	fmt.Println("\n=== worker pool (3 workers, 9 jobs) ===")
	start := time.Now()
	res := workerPool([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 3)
	fmt.Println(res, "ใช้เวลา", time.Since(start).Round(10*time.Millisecond), "(~3 รอบ × 10ms)")

	fmt.Println("\n=== select ===")
	demoSelect()

	fmt.Println("\n=== mutex ===")
	v := Visits{count: map[string]int{}}
	var wg sync.WaitGroup
	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v.Hit("/home")
		}()
	}
	wg.Wait()
	fmt.Println("hits:", v.count["/home"], "(ลองลบ Lock ออกแล้วรัน go run -race ./10-concurrency ดู)")
}
