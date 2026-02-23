module star-tex.org/x/tex

go 1.25.0

retract (
	// v0.3.0 had commits missing "sign-off"
	v0.3.0

	// v0.1.x are the old CGo-based versions
	v0.1.1
	v0.1.0
)

require (
	codeberg.org/go-pdf/fpdf v0.11.1
	git.sr.ht/~sbinet/cmpimg v0.1.0
	golang.org/x/image v0.39.0
	modernc.org/knuth v0.5.5
)

require (
	modernc.org/token v1.1.0 // indirect
	rsc.io/pdf v0.1.1 // indirect
)
