package adapters

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/4strodev/wiring/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	scanner               = bufio.NewScanner(os.Stdin)
	end     chan struct{} = make(chan struct{})
)

func readLine() (string, bool) {
	ok := scanner.Scan()
	if !ok {
		err := scanner.Err()
		if err != nil {
			panic(err)
		}
		return "", false
	}
	return scanner.Text(), true
}

type CliAdapter struct {
	Logger     *slog.Logger
	Connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

// Init implements Adapter.
func (c *CliAdapter) Init(container pkg.Container) error {
	scanner.Split(bufio.ScanLines)
	err := container.Fill(c)
	if err != nil {
		return err
	}
	ch, err := c.Connection.Channel()
	if err != nil {
		return err
	}
	c.channel = ch

	queue, err := c.channel.QueueDeclare("task_queue", true, false, false, false, nil)
	if err != nil {
		return err
	}
	c.queue = queue

	return nil
}

// Start implements Adapter.
func (c *CliAdapter) Start() error {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			fmt.Printf(" > ")
			line, ok := readLine()
			if !ok {
				return
			}
			if line == "exit" {
				break
			}

			err := c.sendMessage(line)
			if err != nil {
				panic(err)
			}
			fmt.Println(line)
		}
	}()

	go func() {
		wg.Wait()
		end <- struct{}{}
	}()

	<-end
	return nil
}

// Stop implements Adapter.
func (c *CliAdapter) Stop() error {
	end <- struct{}{}
	return c.channel.Close()
}

func (c *CliAdapter) sendMessage(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := c.channel.PublishWithContext(ctx, "", c.queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(message),
	})
	if err != nil {
		return err
	}

	c.Logger.Info("message sent")
	return nil
}
