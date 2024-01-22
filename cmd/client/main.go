package main

import (
	"errors"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"slavic-cats-names/internal"
	"slavic-cats-names/internal/task"
	"strconv"
	"syscall"
	"time"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	ticker := time.NewTicker(time.Second)

	for {
		err := request()
		if err != nil {
			log.Println("request error:", err)
		}

		ticker.Reset(5*time.Second + time.Duration(r.Int()%15)*time.Second)
		select {
		case <-c:
			return
		case <-ticker.C:
		}
	}

}

func request() error {
	conn, err := net.Dial("tcp", "server:2000")
	if err != nil {
		return err
	}
	defer conn.Close()

	taskBytes, err := internal.TCPRead(conn, math.MaxUint16)
	if err != nil {
		return errors.Join(err, errors.New("task reading error"))
	}

	task, err := task.ParseTask(taskBytes)
	if err != nil {
		return errors.Join(err, errors.New("task parsing error"))
	}

	nonce := task.Solve()
	err = internal.TCPWrite(conn, []byte(strconv.FormatInt(nonce, 16)))
	if err != nil {
		return errors.Join(err, errors.New("nonce sending error"))
	}

	catName, err := internal.TCPRead(conn, math.MaxUint16)
	if err != nil {
		return errors.Join(err, errors.New("cat name reading error"))

	}

	log.Println("Received msg: ", string(catName))
	return nil
}
