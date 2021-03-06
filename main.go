package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
)

var (
	_zeroWidthReg = regexp.MustCompile(`[\x{200b}]+`)
	_nocharReg    = []*regexp.Regexp{
		regexp.MustCompile(`[\p{Hangul}]+`),  // kr
		regexp.MustCompile(`[\p{Tibetan}]+`), // tibe
		regexp.MustCompile(`[\p{Arabic}]+`),  // arabic
	}
)

func main() {
	defer timeCost(time.Now())
	// var err error
	// txt := strconv.FormatInt(222222, 10)
	// wm := watermark.NewWatermark("./img/uid_mark.png", txt, 32)
	// if wm.Draw(false) != nil {
	// 	fmt.Printf("%v\n", err)
	// 	return
	// }

	// txt := "golang水印大多数"
	// wm := watermark.NewWatermark("./img/mark.png", txt, 32)
	// if wm.Draw(true) != nil {
	// 	fmt.Printf("%v\n", err)
	// 	return
	// }
	//imgprocessing.HTTPPrint()
	//fmt.Printf("%s\n", "create watermark success ...")
	//rpcx.Start()
	//signalHandler()
	// ss := "   "
	// checkTitle(ss)
	// sum(1, 2)
	// sum(1, 2, 3)
	// nums := []int{1, 2, 3, 4, 5}
	// sum(nums...)
	// fmt.Printf("%v\n", 2*time.Second)
	testGorun()
	testGorun1()
}

func timeCost(start time.Time) {
	terminal := time.Since(start)
	fmt.Println(terminal)
}

func sum(nums ...int) {
	fmt.Print(nums, " ")
	total := 0
	for _, num := range nums {
		total += num
	}
	fmt.Println(total)
}

func signalHandler() {
	// signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		println("get a signal")
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			time.Sleep(time.Second * 2)
			println("exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}

func checkTitle(title string) (ct string, ok bool) {
	title = strings.TrimSpace(title)
	println(utf8.RuneCountInString(title))
	if utf8.RuneCountInString(title) > 80 {
		return
	}
	for _, reg := range _nocharReg {
		if reg.MatchString(title) {
			return
		}
	}
	ct = _zeroWidthReg.ReplaceAllString(title, "")
	ok = true
	return
}

func testGorun() {
	slic := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	ch := make([]chan int, 10)
	for i := 0; i < 10; i++ {
		ch[i] = make(chan int)
		go func(c []chan int, s []int, j int) {
			fmt.Printf("input %v\n", s[j])
			c[j] <- j
		}(ch, slic, i)
	}
	for k := range ch {
		fmt.Printf("output %v\n", <-ch[k])
	}
}

func testGorun1() {
	slic := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	ch := make([]chan int, 10)
	for i := 0; i < 10; i++ {
		ch[i] = make(chan int)
		go func(c []chan int, s []int, j int) {
			fmt.Printf("input %v\n", s[j])
			c[j] <- j
		}(ch, slic, i)
	}
	m := 0
	// 带缓冲的channel需要关闭，要不然会死锁
	for k := range ch {
		m++
		fmt.Printf("output %v\n", <-ch[k])
		// if m >= 10 {
		// 	close(ch)
		// }
	}
}
