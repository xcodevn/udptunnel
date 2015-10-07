/* UDP Tunnel - tun.h
 * 2009 - Felipe Astroza
 * Under GPLv2 license (see LICENSE)
 */

#ifndef TUN_H
#define TUN_H

int tun_create();
unsigned int tun_get_packet(int fd, char *buf, unsigned int bufsize);
void tun_put_packet(int fd, char *buf, unsigned int buflen);
void tun_xor(char *input, char *password, int len);

#endif
