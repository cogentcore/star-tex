@x
primitive("jobname",convert,job_name_code);@/
@y
primitive("creationdate",convert,creationdate_code);@/
@!@:creationdate_}{\.{\\creationdate} primitive@>
primitive("jobname",convert,job_name_code);@/
@z

@x
  othercases print_esc("jobname")
@y
  creationdate_code: print_esc("creationdate");
  othercases print_esc("jobname")
@z

@x Just make sure this is defined.  This is certainly not doing the right thing though.
job_name_code: if job_name=0 then open_log_file;
@y
creationdate_code: cur_val:=sys_time;
job_name_code: if job_name=0 then open_log_file;
@z
