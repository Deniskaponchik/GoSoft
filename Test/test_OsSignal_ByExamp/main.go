package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Уведомление о выходе сигнала работает путем отправки значений `os.Signal` в канал.
	// Мы создадим канал для получения этих уведомлений
	// мы также создадим канал, чтобы уведомить нас, когда программа может выйти
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// signal.Notify регистрирует данный канал для получения уведомлений об указанных сигналах.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Эта горутина выполняет блокировку приема сигналов.
	// Когда она получит его, то распечатает его, а затем уведомит программу, что она может быть завершена.
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		fmt.Println("This is my text")
		done <- true
	}()

	// Программа будет ждать здесь, пока не получит ожидаемый сигнал
	// (как указано в приведенной выше процедуре, отправляющей значение в `done`),
	// и затем завершится.
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
