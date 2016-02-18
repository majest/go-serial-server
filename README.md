Go Serial Server
================

Starts a web server on port 8000, listens to the data on /send and sends it to specified serial port.


```
./go-serial-server -h
Usage of ./go-serial-server:
  -boud int
    	Boud rate (default 115200)
  -name string
    	Serial port name (default "/dev/ttys000")
  -port string
    	Server Port number (default "8000")
```
