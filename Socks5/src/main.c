#include "socks.h"
#include <caaa/argparse.h>
#include <caaa/bstrlib.h>
#include <caaa/dbg.h>

struct tagbstring s5_arg_desc = bsStatic("A socks5 server in c");

ArgumentParser *Add_Arguments(int argc, char *argv[]) {
  ArgumentParser *arg =
      Argparse_New_Argument_Parser(&s5_arg_desc, bfromcstr(argv[0]));
  check(arg != NULL, "Could not create argument Parser");
  Argparse_Add_String(arg, "-H", "ip", "0.0.0.0", "The ip to bind to");
  Argparse_Add_Int(arg, "-P", "port", "1080", "The port to bind to");
  return arg;
error:
  return NULL;
}
