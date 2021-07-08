package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

var (
	address = []string{"localhost:9092"}
	topic   = "test"
)

func main() {
	// syncProducer(address)
	asyncProducer1(address)
	// asyncProducer2(address)
}

//同步消息模式
func syncProducer(address []string) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll //等待服务器所有副本都保存成功后的响应
	config.Producer.Return.Successes = true          //是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Errors = true
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return
	}
	defer p.Close()
	srcValue := "sync: this is a message. index=%d"
	for i := 0; i < 10; i++ {
		value := fmt.Sprintf(srcValue, i)
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.ByteEncoder(strconv.Itoa(i)),
			Value: sarama.ByteEncoder(value),
		}
		part, offset, err := p.SendMessage(msg)
		if err != nil {
			log.Printf("send message(%s) err=%s \n", value, err)
		} else {
			fmt.Fprintf(os.Stdout, value+"发送成功，partition=%d, offset=%d \n", part, offset)
		}
		time.Sleep(2 * time.Second)
	}
}

//异步消费者(Goroutines)：用不同的goroutine异步读取Successes和Errors channel
func asyncProducer1(address []string) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	//config.Producer.Partitioner = 默认为message的hash//随机的分区类型
	p, err := sarama.NewAsyncProducer(address, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer error(%v)\n", err)
		return
	}

	//Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var (
		wg                          sync.WaitGroup
		enqueued, successes, errors int
	)

	wg.Add(2) //2 goroutine
	// 发送成功message计数
	go func() {
		defer wg.Done()
		for range p.Successes() {
			successes++
		}
	}()
	// 发送失败计数
	go func() {
		defer wg.Done()
		for err := range p.Errors() {
			log.Printf("%+v 发送失败，err：%s\n", err.Msg, err.Err)
			errors++
		}
	}()
	// 循环发送信息
	// asrcValue := "async-goroutine: this is a message. index=%d"
	var i int
Loop:
	for {
		i++
		// value := fmt.Sprintf(asrcValue, i)
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.ByteEncoder(strconv.Itoa(i)),
			Value: sarama.ByteEncoder(strconv.Itoa(i)),
		}
		select {
		case p.Input() <- msg: // 发送消息
			enqueued++
			fmt.Fprintln(os.Stdout, i)
		case <-signals: // 中断信号
			p.AsyncClose()
			break Loop
		}
		time.Sleep(2 * time.Second)
	}
	wg.Wait()
	fmt.Fprintf(os.Stdout, "发送数=%d，发送成功数=%d，发送失败数=%d \n", enqueued, successes, errors)
}

//异步消费者(Select)：同一线程内，通过select同时发送消息 和 处理errors计数。
//该方式效率较低，如果有大量消息发送， 很容易导致success和errors的case无法执行，从而阻塞一定时间。
//当然可以通过设置config.Producer.Return.Successes=false;config.Producer.Return.Errors=false来解决
func asyncProducer2(address []string) {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	p, err := sarama.NewAsyncProducer(address, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return
	}

	//Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var (
		enqueued, successes, errors, i int
	)
	asrcValue := "async-select: this is a message. index=%d"
Loop:
	for {
		i++
		value := fmt.Sprintf(asrcValue, i)
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.ByteEncoder(strconv.Itoa(i)),
			Value: sarama.ByteEncoder(value),
		}
		select {
		case p.Input() <- msg:
			fmt.Fprintln(os.Stdout, value)
			enqueued++
		case <-p.Successes():
			successes++
		case err := <-p.Errors():
			log.Printf("%+v 发送失败，err：%s\n", err.Msg, err.Err)
			errors++
		case <-signals:
			p.AsyncClose()
			break Loop
		}
		time.Sleep(2 * time.Second)
	}
	fmt.Fprintf(os.Stdout, "发送数=%d，发送失败数=%d \n", enqueued, errors)
}
