@x
@!@:badness_}{\.{\\badness} primitive@>
@y
@!@:badness_}{\.{\\badness} primitive@>
primitive("shellescape",last_item,shellescape_code);@/
@!@:shellescape_}{\.{\\shellescape} primitive@>
@z

@x
  othercases print_esc("badness")
@y
  shellescape_code: print_esc("shellescape");
  othercases print_esc("badness")
@z

@x Read-only status indicating if shell escape is allowed. Expands to zero for off.
  badness_code: cur_val:=last_badness;
@y
  badness_code: cur_val:=last_badness;
  shellescape_code: cur_val := 0;
@z
