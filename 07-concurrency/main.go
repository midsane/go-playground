package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

/*
mutex is used to lock a critical section-> ensure any type of operation-
read/write -> only one can happen at a time.

every operation happens in a order.
for data structuers where both read and write happens very often
*/

// type MidCache struct {
// 	mu   sync.Mutex
// 	data map[int]int
// }

// func (m *MidCache) Get(key int) (int, error) {
// 	m.mu.Lock()
// 	val, ok := m.data[key]
// 	m.mu.Unlock()
// 	var err error
// 	if !ok {
// 		err = errors.New("key does not exist!")
// 	}
// 	return val, err
// }

// func (m *MidCache) Set(key int, val int) {
// 	m.mu.Lock()
// 	m.data[key] = val
// 	m.mu.Unlock()
// }

/*
for data structure where read operation are 90%, write operation is very less
frequent.
we can use readwrite mutex, where simulataneous read operations are allowed.
all write operations are blocked if there are current ongoing read operation, and all read operation +
write operation are blocked when a write is going on
*/
type MidCache struct {
	mu   sync.RWMutex
	data map[int]int
}

func (m *MidCache) Get(key int) (int, error) {
	m.mu.RLock()
	val, ok := m.data[key]
	m.mu.RUnlock()
	var err error
	if !ok {
		err = errors.New("key does not exist!")
	}
	return val, err
}

func (m *MidCache) Set(key int, val int) {
	m.mu.RLock()
	m.data[key] = val
	m.mu.RUnlock()
}

/*
i am assuming all 10 goroutines reports to channel - no panic scenario, if even one
panics, channels wont recieve 10 values, hence deadlock forever
// */
// func main() {
// 	mc := MidCache{data: make(map[int]int)}
// 	dta := make(chan int)
// 	for i := range 10 {
// 		fmt.Println("i: ", i)
// 		go func(i int) {
// 			time.Sleep(time.Second * 1)
// 			mc.Set(1, i)
// 			val, err := mc.Get(1)
// 			if err != nil {
// 				fmt.Print(err)
// 			}
// 			dta <- val
// 		}(i)
// 	}

// 	for range 10 {
// 		fmt.Println(<-dta)
// 	}
// }

/*
a better approach is to range over channel until it keeps receiving value. and explicitly
close the channel once all go routines are done processing via waitgroup
*/

func ImplementCache() {
	mc := MidCache{data: make(map[int]int)}
	dta := make(chan int)

	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Second * 2)
			mc.Set(1, i)
			if val, err := mc.Get(1); err == nil {
				dta <- val
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(dta)
	}()

	for d := range dta {
		fmt.Println(d)
	}
}

/*lets try to simulate a mutex using channel(semaphore style)*/

func Semaphore() {
	semaphore := make(chan struct{}, 1)
	var wg sync.WaitGroup
	counter := 0
	for range 10 {
		//acquiring lock, lock wont be acquired if it still ain't received
		//since its unbuffered
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			//without locking it will cause data race -> can be checked via go run -race .
			semaphore <- struct{}{}
			counter++
			<-semaphore
		}()
	}
	wg.Wait()
	fmt.Println(counter)
}

func main() {
	// SendingWhileClosedLeadsToPanic()
	// WaitOnTwoChannels()
	// Timeout()
	// NonBlockingChannel()
	// SignalInit()
	// BroadCast()
	FanInSimulation()
}

// use cases of unbuffered channels
func WorkerConsumer1() {
	ch := make(chan int)
	go func() {
		jobID := 2
		fmt.Println("worker doing job: ", jobID)
		time.Sleep(time.Second * 3)
		ch <- jobID
	}()

	fmt.Println("job with id:", <-ch, "completed")
}

func WorkerConsumer2() {
	ch := make(chan int)
	go func() {
		for i := range 5 {
			fmt.Println("aded job with id:", i)
			ch <- i
		}
		close(ch)
	}()

	for range 5 {
		fmt.Println("received job of id:", <-ch)
	}

}

/*
can only send to a open channel which is not filled.
in case of unbuffered -> there is no storage, when we send , receiver should be ready
at that exact moment, or else thread will be freezed until we find a reciver.

sending to a closed channel panics.

recieving from a close channel will receive unreceived data first (in case of buffered channel),
after that zero value will be received and value of second param will be false
value, ok := <-ch
ok will be false if channel is close , and value will be 0 value

what if the channel is closed in some go routine while some other go routine is trying to send
via that  channel -> it will make that go routine panic.

example below
*/

func SendingWhileClosedLeadsToPanic() {
	ch := make(chan int)
	fmt.Println("go routine starts")
	go func() {
		time.Sleep(time.Millisecond)
		close(ch)
	}()

	go func() {
		time.Sleep(time.Millisecond * 10)
		fmt.Println(<-ch)
	}()

	ch <- 1

	//since close(ch) will happen before recive, the main thread will panice on ch <- 1
}

// the go routine worker will defintely panics, but the main might terminate that panic surface
// and you might simply see zero value printed
// to see panice change time.Sleep(time.Millisecond) to time.Sleep(time.Second)
func PanicScenario2() {
	ch := make(chan int)
	go func() {
		ch <- 1
	}()

	time.Sleep(time.Millisecond)
	close(ch)

	fmt.Println(<-ch)
}

/*
Select lets a go routine wait on multiple channel at once.
it blocks until one of them is ready.
useful:
	for fan in /fan out architecture
	worker pools
	handling shutdown signals
	timeouts/cancellation
	wait on multiple inputs
*/

//waiting on 2 channels

func WaitOnTwoChannels() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		ch1 <- "from ch1"
	}()

	go func() {
		time.Sleep(time.Second * 2)
		ch2 <- "from ch2"
	}()

	select {
	case msg := <-ch1:
		fmt.Println(msg)
	case msg := <-ch2:
		fmt.Println(msg)
	}

}

// how to implement timeout using select
func Timeout() {
	ch := make(chan string)
	go func() {
		time.Sleep(3 * time.Second)
		ch <- "done"
	}()

	select {
	case msg := <-ch:
		fmt.Println(msg)
	case <-time.After(time.Second * 1):
		fmt.Println("Timeout")
	}
}

//non blocking channel
//by providing default in select-> so that it runs immediately no other cases are available

func NonBlockingChannel() {
	ch := make(chan struct{})

	go func() {
		ch <- struct{}{}
	}()
	time.Sleep(time.Second * 5)
	//remove the above line-> and not ready will be printed instantly.
	select {
	case <-ch:
		fmt.Println("received")
	default:
		fmt.Println("not ready")
	}
}

//struct{}{} -> empty struct of size 0 -> has 0 memory footprint -> just exist a typed symbol in the code.
//therefore used in signaling

//infinite event loop
//lets you handle multiple inputs and exit cleanly

func SignalInit() {
	ch := make(chan int)
	quit := make(chan struct{})

	go func() {
		Signalling(ch, quit)
	}()

	go func() {
		time.Sleep(time.Second * 2)
		quit <- struct{}{}
	}()

	time.Sleep(time.Second * 3)
	ch <- 105
}

func Signalling(ch <-chan int, quit <-chan struct{}) {
	select {
	case msg := <-ch:
		fmt.Println("msg :", msg)
	case <-quit:
		fmt.Println("closing channel")
		//now run a go routine to recivve value from ch so that go routine sending value to it
		//dont free
		//we are do,ing this in a goroutine to ensure non blocking exp
		go func() {
			<-ch
		}()
	}
}

// select blocks until one case is ready
// If multiple ready → random choice
// default makes it non-blocking️
// Closed channels are always ready to receive

// broadcast shutdown
func BroadCast() {
	quit := make(chan struct{})

	go Worker(quit, "worker 1 have completed the task")
	go Worker(quit, "worker 2 have completed the task")
	go Worker(quit, "worker 3 have completed the task")

	time.Sleep(time.Second * 1)
	close(quit)
	//all workers can be stopped at once using close, since all of them will recive a zero value
	//immediately from this channel
	//if i had done quite<-struct{}{} -> then one out of 3 would have stopped only.
	time.Sleep(time.Second * 3)
}

func Worker(quit chan struct{}, str string) {
	fmt.Println("working...")
	msg := make(chan string)
	go func() {
		time.Sleep(time.Second * 2)
		msg <- "work completed"
	}()
	select {
	case m := <-msg:
		fmt.Println(m, "\n", str)
	case <-quit:
		fmt.Println("Stopping")
		go func() {
			<-msg
		}()
		return
	}
}

// basic load balancing, job spread out between workers
func worker(id int, jobs <-chan int) {
	for job := range jobs {
		fmt.Println("worker", id, "processing", job)
		time.Sleep(time.Second)
	}
}

func LoadBalancer() {
	jobs := make(chan int)

	for i := 0; i < 3; i++ {
		go worker(i, jobs)
	}

	for j := 0; j < 6; j++ {
		jobs <- j
	}

	close(jobs)
	time.Sleep(5 * time.Second)
}

// Fan-in: merging multiple channels into one
// receiving/sending to/from a closed channel blocks forever.
// closing a
func FanIn(ch1 <-chan int, ch2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {

		defer func() {
			close(out)
			//after everything must close out, so that range over out works.
		}()

		for ch1 != nil || ch2 != nil {
			select {
			case v, ok := <-ch1:
				if !ok {
					ch1 = nil
					//setting to nil ensure this case will dsiabled now otherwise it will infinetly keep receiving 
					//zero values. receiving and sending from/to nil channel is always blocked. 
					continue
				}
				out <- v

			case v, ok := <-ch2:
				if !ok {
					ch2 = nil
					continue
				}
				out <- v
			}
		}
	}()
	return out
}

func FanInSimulation() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	out := FanIn(ch1, ch2)

	go func() {
		for i := range 5 {
			ch1 <- i
		}
		//closing channel, to make sure false is passed in ok when values are being recevined via select from this channel
		close(ch1)
	}()

	go func() {
		for i := range 6 {
			ch2 <- i
		}
		close(ch2)
	}()

	for v := range out {
		fmt.Println(v)
	}

}
