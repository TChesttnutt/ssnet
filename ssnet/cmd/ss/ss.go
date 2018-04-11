package main

import (
	"net"
	"os"
	"fmt"
	"log"
	"strconv"
	"strings"
	"ssnet/ssnet/ssnet"
	"net/http"
	"io/ioutil"
)

var (
	port     uint
	listener net.Listener
	url string
	filename string
	awgetFlag bool = true
)

func handleConnection(clientJobs chan net.Conn){

	// Wait for the next job to come off the queue.
	back := <- clientJobs

	//States for infomation of stepping stone
	fmt.Println("Added new connection!")

	//makes buffer for stepping stone and for the new chainfile
	var (
	sstones ssnet.SteppingStones
	buf = make([]byte, 10000)
	)

	// will listen for message to process ending in newline (\n)
	_ , err:= back.Read(buf)
	if err != nil {
		log.Fatal("Read: ", err)
	}

	//Turns into arguments then into channels
	arguments := strings.Split(string(buf), "\n")
	jumps, _ := strconv.Atoi(arguments[1])

	//if there is still some jumps left it send it randomly to the next ss node
	if ( jumps != 0){
		for i := 2; i < (jumps+2); i++ {
			stoneInfo := strings.Split(arguments[i], ":")
			fmt.Println("Stones: ", stoneInfo)
			if stoneInfo == nil {
				log.Fatal("StoneInfo: ", err)
			}
			ss, err := ssnet.NewSteppingStone(stoneInfo[0], stoneInfo[1])
			if err != nil {
				log.Fatal("NewSteppingStone: ", err)
			}
			sstones = append(sstones, ss)
		}

		fmt.Println("\tRequest: ", arguments[0])
		fmt.Println("\tChainlist is:")
		for i := 2; i < (jumps+2); i++ {
			fmt.Println("\t", arguments[i])
			}



		//sends to the next Stepping Stone
		err := ssnet.SendReqToRandomStone(arguments[0],sstones, back)
		if err != nil {
			log.Fatal("SendReqToRandomStone: ", err)
		}


	} else{
			//Intro if there is no more jumps left
			fmt.Println("\tRequest: ", arguments[0])
			fmt.Println("\tChainlist is empty!")

			filename, _ := ssnet.GetfileName(arguments[0])

			// if there is no more jumps perform get
			fmt.Println("\tIssuing wget for file <", filename ,">")
			resp, err := http.Get(arguments[0])
			if err != nil {
				log.Fatal("http.Get: ", err)
			}
			defer resp.Body.Close()

			//grabs the body of the file specified in the url
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("ioutil.ReadAll: ", err)
			}
			fmt.Println("..")
			fmt.Println("\tFile Recieved!")
			bodyLength := strconv.Itoa(len(body))

			//Writes it to the connection before it. Attaches the length of the html get file
			fmt.Println("\tRelaying File...")
			fmt.Fprintf(back, bodyLength+"\n")
			back.Write(body)

			fmt.Println("\tGoodbye!")
			back.Close()
		}
}

func main() {

	//sets up channel for incomming Jobs
	clientJobs := make(chan net.Conn)

	//gets Hostname of server
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	//Takes in Port Arguments
	argsProg := os.Args
	if (len(argsProg) != 2){
		log.Fatal("Your arguments are not correct: port#")
	}
 	port := argsProg[1]

 	//Starts server using the inserted port
	fmt.Println("Starting SS server on...")
	fmt.Println("Hostname: ", hostname, " Listening on Port: ", port)

	listener, err := net.Listen("tcp", ":" + port)
	if err != nil {
		fmt.Println("Net.Listen: ", err)
	}
	defer listener.Close()

	//For each accepted connection it passes it to the function to handle it
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Net.Accept: ", err)
		}
		//Passes job onto the handle C
		go handleConnection(clientJobs)
		clientJobs <- conn
	}

}
