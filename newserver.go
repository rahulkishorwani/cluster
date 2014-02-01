package cluster

import (
	"encoding/xml"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"time"
	//"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	//"os/exec"
	//"syscall"
)

func New(id string, fnm string) Servermainstruct {

	//iid,ip,port,serv,ownindex:=getithserveraddr(id,fnm)
	iid, _, port, serv, ownindex := getithserveraddr(id, fnm)

	//fmt.Printf("%s %s\n",ip,port)

	inchan := make(chan *Envelope, 20)
	outchan := make(chan *Envelope, 20)

	//inchan:=make(chan *Envelope)
	//outchan:=make(chan *Envelope)

	newservsock, _ := zmq.NewSocket(zmq.DEALER)

	addr := "tcp://*:" + port
	newservsock.Bind(addr)
	//fmt.Printf("Listening on %s : %s\n",ip,port)	
	var newpeersock [10]*zmq.Socket

	time.Sleep(30000 * time.Millisecond)

	for k, _ := range serv {
		if k != ownindex {
			n, _ := zmq.NewSocket(zmq.DEALER)

			str := "tcp://" + serv[k].Ip + ":" + serv[k].Port

			n.Connect(str)

			newpeersock[k] = n
			//fmt.Printf("Connected on %s %s\n",serv[k].Port,str)	

		}

	}
	limit := len(serv)
	mynewpeersock := newpeersock[:limit]

	m := Servermainstruct{servsock: newservsock, peersock: mynewpeersock, pid: iid, peers: serv, ownindexinpeers: ownindex, in: inchan, out: outchan}

	return m
}
func serverfunc(id string, noofmessagestosend int, numberofservers int, wg *sync.WaitGroup) {
	ffnm := "serverlist" + id + ".xml"
	myservermainstruct := New(id, ffnm)

	go myservermainstruct.send(noofmessagestosend, numberofservers)
	go myservermainstruct.receive(noofmessagestosend, numberofservers)

	go myservermainstruct.Sendtooutbox(id, noofmessagestosend, numberofservers)

	var connarr = make([]int, numberofservers)
	for k := 0; k < numberofservers; k++ {
		connarr[k] = 0
	}
	//connarr[0]=0
	//connarr[1]=0
	//connarr[2]=0
	//connarr[3]=0
	//connarr[4]=0

	fnm := "recvd" + id

	tt, _ := strconv.Atoi(id)
	connarr[tt] = 1
	f2, _ := os.OpenFile(fnm, os.O_WRONLY|os.O_CREATE, 0777)
	for i := 0; ; {

		envelope := <-myservermainstruct.Inbox()
		//fmt.Printf("Received msg from %s %d: '%s'\n", envelope.Pid,envelope.MsgId,envelope.Msg)
		fmt.Fprintf(f2, "%s,%s,%d,%s\n", envelope.Pid, id, envelope.MsgId, envelope.Msg)

		if strings.EqualFold(envelope.Msg.(string), "FIN") {
			tmp, _ := strconv.Atoi(envelope.Pid)
			connarr[tmp] = 1
		}

		j := 0
		for j = 0; j < numberofservers; j++ {
			if connarr[j] == 0 {
				break
			}
		}
		if j == numberofservers {
			break
		}
		/*if( connarr[0]==1 && connarr[1]==1 && connarr[2]==1 && connarr[3]==1 && connarr[4]==1 ) {
			break
		}*/

		i++

	}
	f2.Close()

	wg.Done()
	//wg.Wait()	

}

const (
	BROADCAST = "-1"
)

type Envelope struct {
	Pid string

	MsgId int

	Msg interface{}
}

type Myserver interface {
	Pid() string

	Peers() []Server

	Outbox() chan *Envelope

	Inbox() chan *Envelope

	Sendtooutbox(string, int)
}

type Servermainstruct struct {
	servsock        *zmq.Socket
	peersock        []*zmq.Socket
	pid             string
	peers           []Server
	ownindexinpeers int
	in              chan *Envelope
	out             chan *Envelope
}

func (servermainstruct Servermainstruct) Pid() string {
	return servermainstruct.pid
}
func (servermainstruct Servermainstruct) Peers() []Server {
	return servermainstruct.peers
}
func (servermainstruct Servermainstruct) Outbox() chan *Envelope {
	return servermainstruct.out
}
func (servermainstruct Servermainstruct) Inbox() chan *Envelope {
	return servermainstruct.in
}

func (servermainstruct Servermainstruct) Sendtooutbox(id string, noofmessagestosend int, numberofservers int) {

	fnm := "outbuffer" + id
	f2, _ := os.OpenFile(fnm, os.O_WRONLY|os.O_CREATE, 0777)

	for i := 0; i < noofmessagestosend; i++ {
		q := servermainstruct.peers[rand.Intn(numberofservers)].Id
		//fmt.Printf("\nqqq:%s",q)
		if strings.EqualFold(id, q) {
			servermainstruct.Outbox() <- &Envelope{Pid: BROADCAST, MsgId: i, Msg: "BROADCAST"}
			fmt.Fprintf(f2, "%s,%s,%d,%s\n", id, BROADCAST, i, "BROADCAST")

		} else {
			servermainstruct.Outbox() <- &Envelope{Pid: q, MsgId: i, Msg: "hello there"}
			fmt.Fprintf(f2, "%s,%s,%d,%s\n", id, q, i, "hello there")
		}

	}
	servermainstruct.Outbox() <- &Envelope{Pid: BROADCAST, MsgId: noofmessagestosend, Msg: "FIN"}
	f2.Close()

}
func (servermainstruct Servermainstruct) send(noofmessagestosend int, numberofservers int) {
	fnm := "send" + servermainstruct.Pid()
	f2, _ := os.OpenFile(fnm, os.O_WRONLY|os.O_CREATE, 0777)
	for {

		d := <-servermainstruct.out
		tmp := d.Pid
		d.Pid = servermainstruct.Pid()

		b := "Id:" + d.Pid + ",MsgId:" + strconv.Itoa(d.MsgId) + ",Msg:" + d.Msg.(string)

		for k, _ := range servermainstruct.peers {

			if k != servermainstruct.ownindexinpeers {
				if strings.EqualFold(tmp, "-1") {
					//fmt.Printf("Sending to %d",k)
					servermainstruct.peersock[k].Send(b, 0)
					fmt.Fprintf(f2, "%s,%s,%d,%s\n", d.Pid, servermainstruct.peers[k].Id, d.MsgId, d.Msg.(string))

				} else {

					if strings.EqualFold(tmp, servermainstruct.peers[k].Id) {
						//fmt.Printf("Sending to %d",k)						
						servermainstruct.peersock[k].Send(b, 0)
						fmt.Fprintf(f2, "%s,%s,%d,%s\n", d.Pid, servermainstruct.peers[k].Id, d.MsgId, d.Msg.(string))
					}
				}
			}
		}
		if strings.EqualFold(d.Msg.(string), "FIN") {
			break
		}

	}
	f2.Close()
}

func (servermainstruct Servermainstruct) receive(noofmessagestosend int, numberofservers int) {

	fnm := "inbuffer" + servermainstruct.Pid()
	f2, _ := os.OpenFile(fnm, os.O_WRONLY|os.O_CREATE, 0777)

	var connarr = make([]int, numberofservers)
	for k := 0; k < numberofservers; k++ {
		connarr[k] = 0
	}

	for {

		servermainstruct.servsock.SetRcvtimeo(-1)
		msg, err := servermainstruct.servsock.Recv(0)

		if err != nil {
			estr := "Receiver problem in " + servermainstruct.pid
			panic(estr)
		}

		v := strings.Split(msg, ",")
		vv := strings.Split(v[0], ":")
		vvv := strings.Split(v[1], ":")
		vvvv := strings.Split(v[2], ":")
		o, _ := strconv.Atoi(vvv[1])
		q := Envelope{Pid: vv[1], MsgId: o, Msg: vvvv[1]}

		//fmt.Printf("Received:%s",msg)
		fmt.Fprintf(f2, "%s,%s,%d,%s\n", vv[1], servermainstruct.Pid(), o, vvvv[1])
		servermainstruct.Inbox() <- &q

		if strings.EqualFold(vvvv[1], "FIN") {
			tmp, _ := strconv.Atoi(vv[1])
			connarr[tmp] = 1
		}
		j := 0
		for j = 0; j < numberofservers; j++ {
			if connarr[j] == 0 {
				break
			}
		}
		if j == numberofservers {
			break
		}

	}
	f2.Close()
}

type Serverinfo struct {
	XMLName    Servermeta `xml:"serverinfo"`
	Servermeta Servermeta `xml:"servermeta"`
	Serverlist Serverlist `xml:"serverlist"`
}
type Servermeta struct {
	Servercount int `xml:"servercount"`
}
type Serverlist struct {
	Server []Server `xml:"server"`
}
type Server struct {
	Id   string `xml:"id"`
	Ip   string `xml:"ip"`
	Port string `xml:"port"`
}

/*max capacity 10000 words*/
func getithserveraddr(id string, fnm string) (string, string, string, []Server, int) {
	//xmlFile, err := os.Open("serverlist.xml")
	xmlFile, err := os.Open(fnm)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "0", "0", "0", nil, -1
	}
	defer xmlFile.Close()
	data := make([]byte, 10000)
	count, err := xmlFile.Read(data)
	if err != nil {
		fmt.Println("Can't read the data", count, err)
		return "0", "0", "0", nil, -1
	}
	var q Serverinfo
	xml.Unmarshal(data[:count], &q)
	checkError(err)

	for k, sobj := range q.Serverlist.Server {
		if strings.EqualFold(sobj.Id, id) {
			return q.Serverlist.Server[k].Id, q.Serverlist.Server[k].Ip, q.Serverlist.Server[k].Port, q.Serverlist.Server, k

		}
	}

	return "0", "0", "0", nil, -1
}
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
