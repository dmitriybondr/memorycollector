package main
import (
	"fmt"
	"golang.org/x/crypto/ssh"
"github.com/guillermo/go.procmeminfo"
	"log"
	"io/ioutil"
	"time"
	"statsdclient"
)
type Memory struct {
	Host string
	Memory *procmeminfo.MemInfo
}
func main() {
//Construct map of your hosts.
//TODO it is better to set some conf file for this
hostmap:=make(map[string]string)
//add your hosts here
hostmap["1"]="10.1.8.19"
hostmap["2"]="10.1.9.30"
//set username for ssh connection
username:="dkorol"
//set path to id_rsa key (here is mine)
rsaPath:="/Users/dmitriy_korol/.ssh/id_rsa"


statsdClient:=statsdclient.Client{
	Address: "127.0.0.1:8125",
	Timeout: 30,
	Prefix: "infra_hosts_golang",
}
statsdClient.Connect()
// close the  connection properly if you don't need it anymore
defer statsdClient.Close()





	quit := make(chan struct{})
	sshvalue := make(chan *Memory)
	metrics:= make (map[string]*statsdclient.Metric)
// cosntruct metrics map
	for _, host := range hostmap   {
		metrics[host] = &statsdclient.Metric{Name: host, Value: 0 }
	}
// run remote (cat /proc/meminfo) in parallel
	for _, host := range hostmap  {
		go goSSH(host,rsaPath,username,sshvalue,quit)
	}

	for {
		select {
		case data := <-sshvalue:
			metrics[data.Host].Value=data.Memory.Available()/ 1024.0 /1024.0
            fmt.Println(data.Memory.Available()/ 1024.0 /1024.0)
			go statsdClient.Send(*metrics[data.Host])
		case <- quit:
			return
		}
	}


}



func goSSH (host string, rsaPath string, user string, sshvalue chan *Memory, quit  chan struct{},  )  {
	tickerSSH := time.NewTicker(30 * time.Second)
	log.Printf("Processing Host: %s" , host)
	for {
		select {
		case <-tickerSSH.C:
	var hostKey ssh.PublicKey

	key, err := ioutil.ReadFile(rsaPath)
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}
	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig,
	// and provide a HostKeyCallback.
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
	hostPort:=fmt.Sprintf("%s:%s", host, "22")
	client, err := ssh.Dial("tcp", hostPort, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}

	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
		stdout, err := session.StdoutPipe()
		if err := session.Run("cat /proc/meminfo"); err != nil {
			log.Fatal("Failed to run: " + err.Error())
		}
		meminfo := &procmeminfo.MemInfo{}
		meminfo.Update(stdout)
		result:=&Memory{Host:host,Memory:meminfo}
		sshvalue <- result
		case <-quit:
			return
	}

}
}



