package main

//#include<stdio.h>
//void inC() {
//    printf("I am in C code now!\n");
//}
import "C"
import "fmt"

func main() {
	fmt.Println("I am in Go code now!")
	C.inC()
}
