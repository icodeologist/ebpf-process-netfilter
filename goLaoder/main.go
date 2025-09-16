package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

func main() {
	// remove mem limit to ebpf
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal("Failed to remove the mem limit")
	}

	// load the ebpf(obj/compiled file of c) program
	spec, err := ebpf.LoadCollectionSpec("control_myprocess.o")
	if err != nil {
		log.Fatalf("Failed to load the ebpf program %v\n", err.Error())
	}

	// load the collection
	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		log.Fatalf("Failed to create collection %v\n", err.Error())
	}

	defer coll.Close()

	// get the function from the ebpf file
	prog := coll.Programs["filter_myprocess_requests"]
	if prog == nil {
		log.Fatal("Program 'filter_myprocess_requests' not found")
	}

	// now attach the prog to cgroup
	cgroupPath := "/sys/fs/cgroup/"
	if len(os.Args) > 1 {
		cgroupPath = os.Args[1]
	}
	l, err := link.AttachCgroup(link.CgroupOptions{
		Path:    cgroupPath,
		Attach:  ebpf.AttachCGroupInet4Connect,
		Program: prog,
	})
	if err != nil {
		log.Fatalf("Failed to attach to cgroup %v\n", err)
	}
	defer l.Close()

	fmt.Printf(" eBPF program attached to cgroup: %s\n", cgroupPath)
	fmt.Println(" Process 'myprocess' can now only connect to port 4040")

	fmt.Println("Press CTL+C to exit")

	// wait for the interuption

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("Program detached successfully")
}
