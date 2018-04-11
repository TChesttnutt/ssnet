package main

import (
	"flag"
	"fmt"
	"log"
	"ssnet/ssnet/ssnet"
	"net"
	"os"
)

var (
	chainFileName string
	back net.Conn
)

func main() {

	//grabs url and chain file.
	argsProg := os.Args
	if (len(argsProg) != 2) && (len(argsProg) != 3){
		log.Fatal("Your arguments are not correct: URL -c=[chainfile]")
	}

	url := argsProg[1]
	chainFileName := flag.String("c", "ssnet/examples/chaingang.txt","Chainfile with vaild stepping stones.")
	flag.Parse()


	//Sets up awget
	fmt.Println("awget:")
	fmt.Println("\tRequest: ", url)

	//Reads file to grab array of stepping stones
	sstones, err := ssnet.ReadChainFile(*chainFileName)
	if err != nil {
		log.Fatal("ReadChainFile: ",err)
	}

	//Selects a stepping stone at random and sends request
	err = ssnet.SendReqToRandomStone(url,sstones, back)
	if err != nil {
		log.Fatal("SendReqToRandomStone: " ,err)
	}

}
