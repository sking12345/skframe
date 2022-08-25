package routes

import "fmt"

func TCPNewConnect(fd int) {
	fmt.Println("new fd:", fd)
}
func TCPCloseConnect(fd int) {
	fmt.Println("close fd:", fd)
}

func TCPNewMessage(fd int, data []byte) {
	fmt.Println("new msg:", fd)
	fmt.Println(data)
}
