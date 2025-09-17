## eBPF program to control "myprocess" network system.

## What its doing?
- Intercepting all connection attempts from "myprocess"
- Correctly identifying the destination ports
- Allowing only port 4040
- Blocking all other ports (53, 80, etc.)

## Setting up the test server with "myprocess" name.
- Simple go-server which makes connection to real server on different ports.
- To build it simply run
```
    go build -o myprocess tester.go

```

## HOW TO TEST IT
#### 1. Terminal 1 - Start monitoring eBPF logs.
```
    sudo cat /sys/kernel/debug/tracing/trace_pipe
```
#### 2. Compile your ebpf program.
```
    clang -O2 -target bpf -g -c <ebf_file.c> -o <objectfile.o>
```
#### 3. Load your ebpf program with go.
-  its in /goLoader dir
```
    go build -o loader main.go
```
- then load it with sudo permissions
```
    sudo ./loader
```
#### 4. Finally run your myprocess which you built.
*** Please make your own test if needed ***
```
    sudo ./myprocess
```

#### Now you will be able to see the kernel logs in trace_pipe

#### This was super fun and super hard. I never tried to do anything related to networking and I super interested in this. I love eBPF code which helps me do something awesome. I still have alot to learn and explore. :w

