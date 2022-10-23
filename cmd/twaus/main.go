package main

//#include<stdio.h>
//void inC() {
//    printf("I am in C code now!\n");
//}
import "C"
import (
	"fmt"
	"github.com/slim-bean/twaus/pkg/sensors/sen54"
	"log"
	"periph.io/x/conn/v3/i2c/i2creg"
)

func main() {
	fmt.Println("I am in Go code now!")
	C.inC()

	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatalf("Failed while opening bus: %v", err)
	}
	defer bus.Close()

	_ = sen54.New(bus)
}
