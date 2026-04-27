@x
@d eTeX_int=badness_code+1 {first of \eTeX\ codes for integers}
@y
@d shellescape_code = badness_code + 1 {code for \.{\\shellescape}}
@d eTeX_int=shellescape_code+1 {first of \eTeX\ codes for integers}
@z

@x
@d etex_convert_codes=etex_convert_base+1 {end of \eTeX's command codes}
@d job_name_code=etex_convert_codes {command code for \.{\\jobname}}
@y
@d etex_convert_codes=etex_convert_base+1 {end of \eTeX's command codes}
@d expanded_code = etex_convert_codes {command code for \.{\\expanded}}
@d jstex_first_expand_code = expanded_code + 1 {base for \jsTeX's command codes}
@d strcmp_code = jstex_first_expand_code + 0 {command code for \.{\\strcmp}}
@d creationdate_code = jstex_first_expand_code + 1 {command code for \.{\\creationdate}}
@d filesize_code = jstex_first_expand_code + 2 {command code for \.{\\filesize}}
@d jstex_convert_codes = jstex_first_expand_code + 3 {end of \jsTeX's command codes}
@d job_name_code=jstex_convert_codes {command code for \.{\\jobname}}
@z
