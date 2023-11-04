package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	localMultiAddr, err := multiaddr.NewMultiaddr(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port),
	)
	host, err := libp2p.New(
		libp2p.ListenAddrs(localMultiAddr),
	)

	host.SetStreamHandler("/echo/1.0.0", handleStream)

	if err != nil {
		logrus.WithError(err).Errorln("Create libp2p host failed.")
		return
	}

	logrus.Infoln("Create new libp2p host, listen addresses: ", host.Addrs())
	logrus.Infof("Node address: /ip4/127.0.0.1/tcp/%v/p2p/%s", port, host.ID().String())

	if peerAddr != "" {
		ctx := context.Background()

		maddr, err := multiaddr.NewMultiaddr(peerAddr)
		if err != nil {
			logrus.WithError(err).Errorln("Convert address from string failed.")
			return
		}

		addr, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			logrus.WithError(err).Errorln("Get address info from multiple address failed.")
			return
		}
		err = host.Connect(ctx, *addr)
		if err != nil {
			logrus.WithError(err).Errorln("Connect to new node failed.")
			return
		}

		logrus.Infof("Connect to node %s success.", peerAddr)

		stream, err := host.NewStream(ctx, addr.ID, "/echo/1.0.0")
		if err != nil {
			logrus.WithError(err).Errorln("Create stream failed.")
			return
		}
		go sendMessage(stream)
	}

	c := make(chan os.Signal)
	select {
	case sign := <-c:
		logrus.Infof("Got %s signal. Aborting...", sign)

		if err := host.Close(); err != nil {
			logrus.WithError(err).Errorln("Close host errored.")
		}
	}
}

func sendMessage(s network.Stream) {
	var msg string
	buf := bufio.NewReader(s)
	for {
		_, err := fmt.Scan(&msg)
		if err != nil {
			logrus.WithError(err).Errorln("Read data from console failed.")
			break
		}
		msg += "\n"

		_, err = s.Write([]byte(msg))
		if err != nil {
			logrus.WithError(err).Errorln("Write message to stream failed.")
			break
		}

		str, err := buf.ReadString('\n')
		logrus.Infof("Read from stream: %s", str[:len(str)-1])
	}
}

func handleStream(s network.Stream) {
	for {
		buf := bufio.NewReader(s)
		str, err := buf.ReadString('\n')
		if err != nil {
			logrus.WithError(err).Errorln("Receive failed, stream routine exit.")
			break
		}

		logrus.Infof("Read from stream: %s", str[:len(str)-1])
		_, err = s.Write([]byte(str))
		if err != nil {
			logrus.WithError(err).Errorln("Write to stream failed, routine exit.")
		}
	}
}
