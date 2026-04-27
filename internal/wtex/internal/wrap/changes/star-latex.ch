@x
@d eTeX_version=2 { \.{\\eTeXversion} }
@y
@d StarTeX_version=0 { \.{\\StarTeXversion} }
@d StarTeX_revision==".8.x" { \.{\\StarTeXrevision} }
@d StarTeX_version_string=='-0.8.x' {current \StarTeX\ version}
@d eTeX_version=2 { \.{\\eTeXversion} }
@z

@x
@d eTeX_banner=='This is e-TeX, Version 3.141592653',eTeX_version_string
  {printed when \eTeX\ starts}
@#
@d TeX_banner=='This is TeX, Version 3.141592653' {printed when \TeX\ starts}
@#
@d banner==eTeX_banner
@#
@d TEX==ETEX {change program name into |ETEX|}
@y
@d StarTeX_banner=='This is Star-TeX, Version 3.141592653',StarTeX_version_string
  {printed when \StarTeX\ starts}
@#
@d TeX_banner=='This is TeX, Version 3.141592653' {printed when \TeX\ starts}
@#
@d banner==StarTeX_banner
@#
@d TEX==STARTEX {change program name into |STARTEX|}
@z

@x [2.23] l.723 - Translate characters if desired, otherwise allow them all.
{Initialize |xchr| to the identity mapping.}
for i:=0 to @'37 do xchr[i]:=i;
for i:=@'177 to @'377 do xchr[i]:=i;
@y
{Initialize |xchr| to the identity mapping.}
for i:=0 to @'37 do xchr[i]:=chr(i);
for i:=@'177 to @'377 do xchr[i]:=chr(i);
@z
