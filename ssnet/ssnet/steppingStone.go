package ssnet

import (
	"net"
	"fmt"
	"os"
	"bufio"
	"log"
	"strings"
	"strconv"
	"math/rand"
	"time"
)

var (
	sstones SteppingStones
	lines []string
	selectedStoneIndex int
	url string
)

type steppingStone struct {
	Addr net.IP
	Port uint16
}

func (s steppingStone) String() string {
	return fmt.Sprintf("%s:%d", s.Addr, s.Port)
}

func SendReqToRandomStone(u string, sstones SteppingStones, back net.Conn) ( err error) {

	//sets the global url for this connection
	url =u

	//sees if the list is empty
	if len(sstones) != 0 {
		//grabs seed so random is not deterministic
		rand.Seed(time.Now().UnixNano())
		selectedStoneIndex := rand.Intn(len(sstones))
		fmt.Println("\tNext SS is ", sstones[selectedStoneIndex])

		//Connects to the given ss server if it is set up
		forward, err := net.Dial("tcp", sstones[selectedStoneIndex].String())
		if err != nil {
			log.Fatal("Net.Dial: ", err)
		}

		//sends the number of servers left to the socket
		newLength := len(sstones) - 1
		newLengthS := strconv.Itoa(newLength)
		fmt.Fprintf(forward, u+"\n")
		fmt.Fprintf(forward, newLengthS+"\n")
		for index, _ := range sstones {
			if (index != selectedStoneIndex) {
				forward.Write([]byte(sstones[index].String() + "\n"))
				if err != nil {
					log.Fatal("conn.Write: ", err)
				}
			}
		}
		if back != nil {
			HandleReqFromConn(back, forward)

		} else{
			cleanUp(forward)
			forward.Close()
		}
	}
	return err
}

func HandleReqFromConn(back net.Conn, forward net.Conn) ( err error) {

	reader := bufio.NewReader(forward)

	//Breaking up the the file
	// will listen for message to process ending in newline (\n) to get size
		fmt.Println("\tWaiting for file...")
		length, err :=  reader.ReadString('\n')
		if err != nil {
			log.Fatal("bufio.NewReader-ReadString: ", err)
		}
		// output message received
		length = strings.Trim(length, "\n")
		size, err := strconv.ParseUint(length,0, 64)
		if err != nil {
			log.Fatal("strconv.Atoi: ", err)
		}
		buf := make([]byte, size)

		//Reads the get file from the previous
		_ ,err = forward.Read(buf)
		if err != nil {
			log.Fatal("Read: ", err)
		}

		fmt.Println("..")
		fmt.Println("\tRelaying File...")
		fmt.Fprintf(back,  length + "\n")
		back.Write(buf)

		//Closes the the back File
		fmt.Println("\tGoodbye!")
		back.Close()


	return err
}

type SteppingStones []steppingStone

func NewSteppingStone(ipstr, portstr string) (ss steppingStone, err error) {

		//flips strings into the right types to make a stepping stone
		ip := net.ParseIP(ipstr)
		port, err := strconv.Atoi(portstr)
		if err != nil {
			log.Fatal("strconv.Atoi(portstr): ", err)
		}

		ss = steppingStone{ip, uint16(port)}

		return ss, err
}

func ReadChainFile(fname string) (SteppingStones, error) {

	//Opens file
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal("File could not be found",err)
	}
	defer file.Close()

	//reads file and places it in chainFile reader buffer
	reader := bufio.NewReader(file)
	if err != nil {
		log.Fatal("bufio.NewReader:", err)
	}

	fmt.Println("\tChainlist is:",)
	//Reads chain file line by line using buffer and puts it into the array lines
	line, isPrefix ,err := reader.ReadLine()
	for err == nil && !isPrefix {
		lines = append(lines, string(line))
		line, isPrefix , err = reader.ReadLine()
		if (strings.Compare(string(line),"\n") != 0){
		fmt.Println("\t",string(line))}
	}

	//Splits the ips and ports to create stepping stones
	length, err:= (strconv.Atoi(lines[0]))
	for i:=1; i<= length; i++{
		stoneInfo := strings.Fields(lines[i])
		ss, err := NewSteppingStone(stoneInfo[0], stoneInfo[1])
		if err != nil {
			log.Fatal("newSteppingStone: ", err)
		}
		sstones = append(sstones, ss)
	}
	return sstones ,err
}
func cleanUp( forward net.Conn){

	fmt.Println("\tWaiting for file...")
	reader := bufio.NewReader(forward)

	//Breaking up the the file
	// will listen for message to process ending in newline (\n) to get size
	length, err :=  reader.ReadString('\n')
	if err != nil {
		log.Fatal("bufio.NewReader-ReadString: ", err)
	}
	// output message received
	length = strings.Trim(length, "\n")
	size, err := strconv.ParseUint(length,0, 64)
	if err != nil {
		log.Fatal("strconv.Atoi: ", err)
	}
	buf := make([]byte,  size)

	//Reads the get file from the previous
	_ ,err = forward.Read(buf)
	if err != nil {
		log.Fatal("Read: ", err)
	}

	fmt.Println("..")

	//Grabs the file name
	fileName, _ := GetfileName(url)

	//Creates file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Create File: ", err)
	}

	// Writes to the actual file
	file.Write(buf)
	fmt.Println("\tReceived file <", fileName,">")
	fmt.Println("\tGoodbye!")
	file.Close()
}

//Names file
func GetfileName (url string) (fileName string, err error) {

	lastSlash := strings.LastIndex(url, "/")
	if (lastSlash < 0) || (lastSlash == 6) ||(lastSlash == 7) {
		fileName = "index.html"
	} else if (len(url[lastSlash:]) <= 1) {
		fileName = "index.html"
		} else {
		fileName = url[lastSlash+1:]
	}

	return fileName, err
}