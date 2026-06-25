@x
@ The procedure |conv_toks| uses |str_toks| to insert the token list
for |convert| functions into the scanner; `\.{\\outer}' control sequences
are allowed to follow `\.{\\string}' and `\.{\\meaning}'.

@p procedure conv_toks;
var old_setting:0..max_selector; {holds |selector| setting}
@!c:number_code..job_name_code; {desired type of conversion}
@!save_scanner_status:small_number; {|scanner_status| upon entry}
@!b:pool_pointer; {base of temporary string}
begin c:=cur_chr; @<Scan the argument for command |c|@>;
@y
@ The procedure |conv_toks| uses |str_toks| to insert the token list
for |convert| functions into the scanner; `\.{\\outer}' control sequences
are allowed to follow `\.{\\string}' and `\.{\\meaning}'.

@d save_cur_string==if str_start[str_ptr]<pool_ptr then u:=make_string else u:=0
@d restore_cur_string==if u<>0 then decr(str_ptr)

@p procedure conv_toks;
var old_setting:0..max_selector; {holds |selector| setting}
@!c:number_code..job_name_code; {desired type of conversion}
@!u: str_number;
@!s: str_number;
@!saved_chr:halfword;
@!save_warning_index, @!save_def_ref:pointer;
@!cat:small_number; {desired catcode, or 0 for automatic |spacer|/|other_char| selection}
@!save_scanner_status:small_number; {|scanner_status| upon entry}
@!b:pool_pointer; {base of temporary string}
begin c:=cur_chr; cat:=0; @<Scan the argument for command |c|@>;
@z

@x
@<Scan and build the body of the token list; |goto found| when finished@>;
found: scanner_status:=normal;
if hash_brace<>0 then store_new_token(hash_brace);
scan_toks:=p;
end;
@y
@<Scan and build the body of the token list; |goto found| when finished@>;
found: scanner_status:=normal;
if hash_brace<>0 then store_new_token(hash_brace);
scan_toks:=p;
end;

@ @d call_func(#) == begin if # <> 0 then do_nothing end
@d flushable(#) == (# = str_ptr - 1)

@p procedure flush_str(s: str_number); {flush a string if possible}
begin
    if flushable(s) then
        flush_string;
end;

function tokens_to_string(p: pointer): str_number; {return a string from tokens
list}
begin
  if selector = new_string then begin
    help1("tokens_to_string() called while selector = new_string"); error;
  end;
  old_setting:=selector; selector:=new_string;
  show_token_list(link(p),null,pool_size-pool_ptr);
  selector:=old_setting;
  tokens_to_string:=make_string;
end;
@z
