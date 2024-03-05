
![alt text](https://github.com/antelman107/net-wait-go/blob/master/tube2.svg?raw=true)

# What is it?
Both utility and GO package to wait for port to open (TCP, UDP).

# Why do we need it?

In dockerized applications usually we deploy 
several containers with the main program container.

We need to know if containers are started, so we
can continue to execute our program or fail it after some deadline.

There are a lot of examples through the internet 
that advise us to use several `bash` commands like `nc`, `timeout`, `curl`. But what if we have GO program minimal docker image `from scratch` that does not have `bash`? We could use this package as library in our program - just put couple lines of code to check if some ports are available.

But also this package can be donwloaded as utility and used from command line.

# Library usage

## Simple
```GO
import "github.com/antelman107/net-wait-go/wait"

if !wait.New().Do([]string{"postgres:5432"}) {
    logger.Error("db is not available")
    return
}
```

## All optional settings definition
```GO
import "github.com/antelman107/net-wait-go/wait"

if !wait.New(
      wait.WithProto("tcp"),
      wait.WithWait(200*time.Millisecond),
      wait.WithBreak(50*time.Millisecond),
      wait.WithDeadline(15*time.Second),
      wait.WithDebug(true),
).Do([]string{"postgres:5432"}) {
    logger.Error("db is not available")
    return
}
```

# Utility usage

```bash
$ go get github.com/antelman107/net-wait-go
$ net-wait-go

  -addrs string
        address:port(,address:port,address:port,...)
  -deadline uint
        deadline in milliseconds (default 10000)
  -debug
        debug messages toggler
  -delay uint
        break between requests in milliseconds (default 50)
  -packet string
        UDP packet to be sent
  -proto string
        tcp (default "tcp")
  -wait uint
        delay of single request in milliseconds (default 100)
```

## 1 service check
```bash
net-wait-go -addrs ya.ru:443 -debug true

2020/06/30 18:07:38 ya.ru:443 is OK

# return code is 0
```

## 2 services check
```bash
net-wait-go -addrs ya.ru:443,yandex.ru:443 -debug true

2020/06/30 18:09:24 yandex.ru:443 is OK
2020/06/30 18:09:24 ya.ru:443 is OK

# return code is 0 (if all services are OK)
```

## 2 services check (fail)
```bash
net-wait-go -addrs ya.ru:445,yandex.ru:445 -debug true

2020/06/30 18:09:24 yandex.ru:445 is FAILED
2020/06/30 18:09:24 ya.ru:445 is is FAILED
...
# return code is 1 (if at least 1 service is failed)
```

# UDP support
Since UDP as protocol does not provide connection between a server and clients,
it is not supported in the most of popular
utilities:
 - `wait-for-it` issue - https://github.com/vishnubob/wait-for-it/issues/29)
 - `netcat` (`nc`) has following note in its manual page:
     ```
    CAVEATS
           UDP port scans will always succeed (i.e. report the port as open)
    ```
   
`net-wait-go` provides UDP support, working following way:
 - sends a meaningful packet to the server
 - waits for a message back from the server (1 byte at least)
 
## UDP library usage example
Counter Strike game server is accessible via UDP.
Let's check random Counter Strike server
 by sending A2S_INFO packet (https://developer.valvesoftware.com/wiki/Server_queries#A2S_INFO)

 ```go
e := wait.New(
      wait.WithProto("udp"),
      wait.WithUDPPacket([]byte{
            0xFF, 0xFF, 0xFF, 0xFF, 0x54, 0x53, 
            0x6F, 0x75, 0x72, 0x63, 0x65, 0x20, 
            0x45, 0x6E, 0x67, 0x69, 0x6E, 0x65, 
            0x20, 0x51, 0x75, 0x65, 0x72, 0x79, 
            0x00}),
      wait.WithDebug(true),
      wait.WithDeadline(time.Second*2),
)
if !e.Do([]string{"46.174.53.245:27015","185.158.113.136:27015"}) {
      logger.Error("udp services are not available")
      return
} 
```

`WithUDPPacket` value here is the base64-encoded A2S_INFO packet, which is documented here - https://github.com/wriley/steamserverinfo/blob/master/steamserverinfo.go#L133


## UDP utility usage example 

```bash
net-wait-go -proto udp -addrs 46.174.53.245:27015,185.158.113.136:27015 -packet '/////1RTb3VyY2UgRW5naW5lIFF1ZXJ5AA==' -debug true
 
2020/07/12 15:13:25 udp 185.158.113.136:27015 is OK
2020/07/12 15:13:25 udp 46.174.53.245:27015 is OK

# return code is 0
```

`-packet` value here is the base64-encoded A2S_INFO packet (see above). Since it is problematic to pass
binary value in command line, base64 encoding for 
`packet` parameter is chosen.





