#ifndef SOCKS_H_
#define SOCKS_H_

#include <caaa/argparse.h>
#include <caaa/dbg.h>

enum {
  VERSION = 5,                 // Socks protocol version
  STATE_HANDSHAKE = 0,         // Connection state: in handshake
  STATE_REQUEST = 1,           // Connection state: connecting
  STATE_ESTABLISHED = 2,       // Connection state: established
  HANDSHAKE_NOAUTH = 0,        // Handshake method - no authentication
  HANDSHAKE_GSSAPI = 1,        // Handshake method - GSSAPI auth
  HANDSHAKE_USERPASS = 2,      // Handshake method - user/password auth
  HANDSHAKE_FAILURE = 0xff,    // Handshake method - failure
  CMD_CONNECT = 1,             // Command: CONNECT
  CMD_BIND = 2,                // Command: BIND
  CMD_UDP_ASSOCIATE = 3,       // Command: UDP ASSOCIATE
  ADDR_TYPE_IPV4 = 1,          // Address type: IPv4
  ADDR_TYPE_DOMAIN = 3,        // Address type: Domain name
  ADDR_TYPE_IPV6 = 4,          // Address type: IPv6
  RESP_SUCCESS = 0,            // Response: success
  RESP_FAILURE = 1,            // Response: failure
  RESP_NOT_ALLOWED = 2,        // Response status
  RESP_NET_UNREACHABLE = 3,    // Response status
  RESP_HOST_UNREACHABLE = 4,   // Response status
  RESP_CONN_REFUSED = 5,       // Response status
  RESP_TTL_EXPIRED = 6,        // Response status
  RESP_CMD_NOT_SUPPORTED = 7,  // Response status
  RESP_ADDR_NOT_SUPPORTED = 8, // Response status
};

ArgumentParser *Add_Arguments(int argc, char *argv);

#endif // SOCKS_H_
