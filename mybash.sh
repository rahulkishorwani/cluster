#!/bin/bash
rm -f *~
rm -f inbuffer send outbuffer recvd
cat inbuffer* > inbuffer
cat send* > send
cat outbuffer* > outbuffer
cat recvd* > recvd

sort -t"," -k1 -k2 -k3 outbuffer -o outbuffer
sort -t"," -k1 -k2 -k3 send -o send
sort -t"," -k1 -k2 -k3 inbuffer -o inbuffer
sort -t"," -k1 -k2 -k3 recvd -o recvd

echo ""
echo "Sent in outbuffer from application(with broadcast)"
wc -l outbuffer

echo ""
echo "Actually sent through zmq(broadcast treated as number of peer to peer messages)"
wc -l send

echo ""
echo "Received in inbuffer through zmq"
wc -l inbuffer

echo ""
echo "Actually received in application"
wc -l recvd
