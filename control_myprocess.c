#include "vmlinux.h"
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

#define TARGET_PROCESS "myprocess"
#define ALLOWED_PORT 4040

SEC("cgroup/connect4")
int filter_myprocess_requests(struct bpf_sock_addr *ctx) {
  char comm[TASK_COMM_LEN];

  // first get the current process name
  if (bpf_get_current_comm(&comm, sizeof(comm)) != 0) {
    return 1;
  }
  // then filter only "myprocess"
  if (__builtin_memcmp(comm, TARGET_PROCESS, sizeof(TARGET_PROCESS) - 1) != 0) {
    return 1;
  }
  // else if the process is "myprocess"
  // then get the destination port
  __u16 dest = bpf_ntohs(ctx->user_port);

  // then simply check if its going to ALLOWED_PORT
  bpf_printk("process is trying to connect DEST value %d", dest);
  if (dest == ALLOWED_PORT) {
    bpf_printk("ALLOWED %s is connecting to port %u", comm, dest);
    return 1; // alllow
  } else {
    bpf_printk("BLOCKED %s is connecting to port %u", comm, dest);
    return 0; // block
  }
}

char __license[] SEC("license") = "GPL";
