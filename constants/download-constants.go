package constants

type Size = int64

const (
	Byte     Size = 1
	KiloByte      = 1000 * Byte
	MegaByte      = 1000 * KiloByte
)

type Status struct {
	Id         string
	Downloaded int64
	Total      int64
	FileName   string
	Error      error
}
