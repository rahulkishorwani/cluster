Cluster
=======

Cluster is a repository which has go files for making the interconnections between number of servers(individual server as one cluster) connected to each other & has test file for testing the validity of operations performed by this package.

Language Used
=====
     Go language of Google
Hierarchy
=====

Libraries Used
=====
System:- strconv,strings,fmt,sync,os,os/exec,syscall,time,math/rand,encoding/xml.

Other:-Cluster,zmq4.

Install
=====
go get github.com/rahulkishorwani/cluster

Run
=====
`go test`
 This code will run the code over the number of tests & shows the message of passing of those tests.
 
How to use
=====
     The new method in this package allows you to create the instance of cluster server. This will be one of the server that will be the part of cluster.
`server_main := New(host_id, size_of_in_chan ,size_of_out_chan,fnm,delay_before_conn)`
     where 
          host_id is the id in the form of integer which is present in the configuration file
          size_of_in_chan is the inbox allowable channel size in the form of integer value
          size_of_out_chan is the outbox allowable channel size in the form of integer value
          delay_before_conn is the delay before the connection between the different servers actually happen. It is in the form of int which we should be the value in milliseconds
     
     Then with this instance you have to create the two processes send & receive.
     
`go my_server_main.send(number_of_servers,wg)`
     Send process senses the outbox channel continuously & when the data comes into this channel this process send the data over the cluster to appropriate server.
     where number_of_servers means the number of servers in the configuration file
           wg is of type *sync.WaitGroup so that the calling process can wait for this process to finish
          
`go server_main.receive(number_of_servers,wg)`
     Receive process is always is need of message from network when any message comes onto the receive interface of socket, this method captures that data & converts it into channel envelope object as we will see in the calling function.
     where number_of_servers means the number of servers in the configuration file
           wg is of type *sync.WaitGroup so that the calling process can wait for this process to finish
           
           
`e:=Envelope{ Pid:pid, MsgId:msgId, Msg:msgs }`
     This program creates the instance of structure "Envelope" which contains three things
     Pid : integer storing the process id of the receiving process
     MsgId : It is the message ID of the message
     Msg : This is the interface which is prepared to store the byte array
     
     
`server_main.Outbox()<-e`
     This program sends the Envelope object e over the network

`e:=<-server_main.Inbox()`
     This code blocks until the server gets atleast one message in the inbox. And stores the Envelope object into e.
     

Constraints
=====
   Program have been constrained by the configuration(servelist) xml files. For current server configurations, it works 
for 5 servers. If we want to change the the number of servers then we have to change each xml file to accomodate those much of server entries & you must have to give a corresponding parameter for running the test cases.


Rights
=====
   This code is open source. So anyone can use this.
