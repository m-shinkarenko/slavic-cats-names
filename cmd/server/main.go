package main

import (
	"errors"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"slavic-cats-names/internal"
	gen "slavic-cats-names/internal/slavic-cats-names-gen"
	"slavic-cats-names/internal/task"
)

const (
	maxContentLength = math.MaxUint16
	deadline         = 5 * time.Second
	difficulty       = 20
)

func main() {
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-c
		l.Close()
	}()

	wg := sync.WaitGroup{}
	for {
		conn, err := l.Accept()
		if errors.Is(err, net.ErrClosed) {
			break
		}
		if err != nil {
			log.Println("connection error: ", err)
			continue
		}

		wg.Add(1)
		go func(c net.Conn) {
			defer func() {
				// You don't need to close the connection after each request here, you can reuse it.
				// However, I decided not to complicate things for this task
				c.Close()
				wg.Done()
			}()

			c.SetDeadline(time.Now().Add(deadline))
			if err := handle(conn); err != nil {
				// You can make the response more complex and add something like a status code,
				// but I decided it's not necessary
				internal.TCPWrite(c, []byte(err.Error()))
				log.Println("handle connection error:", err)
			}
		}(conn)
	}

	// Here, you can create a ticker that will forcibly terminate everything
	// if the requests take too long to complete.
	// Again, I decided not to complicate things
	wg.Wait()
}

func handle(c net.Conn) error {
	task, err := task.NewTask(difficulty)
	if err != nil {
		return errors.Join(err, errors.New("task creation error"))
	}

	err = internal.TCPWrite(c, task.ToBytes())
	if err != nil {
		return errors.Join(err, errors.New("task sending error"))
	}

	nonceBytes, err := internal.TCPRead(c, maxContentLength)
	if err != nil {
		return errors.Join(err, errors.New("nonce reading error"))
	}

	nonce, err := strconv.ParseInt(string(nonceBytes), 16, 64)
	if err != nil {
		return errors.Join(err, errors.New("nonce parsing error"))
	}

	if !task.Validate(nonce) {
		return errors.New("task was solved incorrectly")
	}

	err = internal.TCPWrite(c, []byte(gen.GetSlavicCatName()))
	if err != nil {
		return errors.Join(err, errors.New("error writing response"))
	}

	return nil
}
