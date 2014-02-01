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

Constraints
=====
   Program have been constrained by the configuration(servelist) xml files. For current server configurations, it works 
for 5 servers. If we want to change the the number of servers then we have to change each xml file to accomodate those much of server entries & you must have to give a corresponding parameter for running the test cases.

Use
=====
   This code can be useful to build a cluster interface to communicate with other clusters.

Rights
=====
   This code is open source. So anyone can use this.
