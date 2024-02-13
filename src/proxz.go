package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/DataDog/zstd"
)

const bufSize = 1024 * 1024
const zstdLevel = -7

func handleConnection(rdConn, wrConn net.Conn, compress bool) {
	buf := make([]byte, bufSize)

	var rdIo io.ReadCloser
	var wrIo io.WriteCloser

	if compress {
		rdIo = rdConn
		wrIo = zstd.NewWriterLevel(wrConn, zstdLevel)
		defer wrIo.Close()
	} else {
		rdIo = zstd.NewReader(rdConn)
		wrIo = wrConn
		defer rdIo.Close()
	}

	for {
		n, err := rdIo.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		if n > 0 {
			_, err = wrIo.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}

			// Flush the writer or the data may not be sent
			if flusher, ok := wrIo.(interface{ Flush() error }); ok {
				if err := flusher.Flush(); err != nil {
					log.Fatal(err)
				}
			}
		}

		if err == io.EOF {
			break
		}
	}
}

func main() {
	compress := flag.Bool("compress", false, "Enable compression")
	cFlag := flag.Bool("c", false, "Enable compression (shorthand)")
	flag.Parse()

	if *cFlag {
		compress = cFlag
	}

	// Parse non-flag arguments
	args := flag.Args()

	// Handle flags and arguments
	if len(args) < 2 {
		fmt.Println("Usage: proxz [-c|--compress] <(tcp|unix):address> <(tcp|unix):address>")
		os.Exit(1)
	}

	// Split the addresses into network and address
	leftAddr := strings.SplitN(args[0], ":", 2)
	rightAddr := strings.SplitN(args[1], ":", 2)

	// Validate the addresses
	if len(leftAddr) != 2 || len(rightAddr) != 2 {
		fmt.Println("Invalid address")
		os.Exit(1)
	}
	if leftAddr[0] != "tcp" && leftAddr[0] != "unix" {
		fmt.Println("Invalid network")
		os.Exit(1)
	}
	if rightAddr[0] != "tcp" && rightAddr[0] != "unix" {
		fmt.Println("Invalid network")
		os.Exit(1)
	}

	// Create a Unix domain socket file if the address is a Unix domain socket
	if leftAddr[0] == "unix" {
		if err := os.Remove(leftAddr[1]); err != nil && !os.IsNotExist(err) {
			log.Fatal(err)
		}
	}

	// Listen on the left address
	ln, err := net.Listen(leftAddr[0], leftAddr[1])
	if err != nil {
		log.Fatal(err)
	}

	// Create a channel to receive termination signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Run a goroutine that cleans up the socket file when the program is terminated
	go func() {
		<-sigs
		if leftAddr[0] == "unix" {
			os.Remove(leftAddr[1])
		}
		os.Exit(0)
	}()

	// Handle incoming connections
	for {
		// Accept an incoming connection.
		lConn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Connect to the right address
		rConn, err := net.Dial(rightAddr[0], rightAddr[1])
		if err != nil {
			log.Fatal(err)
		}

		// Handle the connection
		go handleConnection(lConn, rConn, *compress)
		go handleConnection(rConn, lConn, !*compress)
	}
}
