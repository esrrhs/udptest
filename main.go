package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"
)

var usage = `
Usage:

    udptest -p 8888

    udptest -p 8888 -t 8899 -s "hello"

`

func main() {
	fmt.Println("start...")

	p := flag.Int("p", 8888, "listen port")
	t := flag.Int("t", 0, "target port")
	sip := flag.String("sip", "127.0.0.1", "server ip")
	s := flag.String("s", "hello", "send string")
	flag.Usage = func() {
		fmt.Printf(usage)
	}

	flag.Parse()

	if flag.NArg() != 0 {
		flag.Usage()
		return
	}

	ipaddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(*p))
	if err != nil {
		fmt.Println(err)
		return
	}

	listener, err := net.ListenUDP("udp", ipaddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	fmt.Printf("listen at %d \n", *p)

	if *t > 0 {

		ipaddrtarget, err := net.ResolveUDPAddr("udp", *sip+":"+strconv.Itoa(*t))
		if err != nil {
			fmt.Println(err)
			return
		}

		data := []byte(*s)
		_, err = listener.WriteToUDP(data, ipaddrtarget)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("send to %v %s \n", ipaddrtarget, *s)
	}

	bytes := make([]byte, 10240)
	for {
		listener.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
		n, srcaddr, err := listener.ReadFromUDP(bytes)
		if err != nil {
			if neterr, ok := err.(*net.OpError); ok {
				if neterr.Timeout() {
					// Read timeout
					continue
				} else {
					fmt.Printf("Error read udp %s\n", err)
					continue
				}
			}
		}

		str := string(bytes[:n])

		fmt.Printf("recv data %s from %s \n", str, srcaddr.String())
	}
}
