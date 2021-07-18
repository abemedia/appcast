package appcast

import "errors"

//go:generate stringer -type=OS -linecomment
type OS uint8

const (
	Unknown   OS = iota //
	MacOS               // macos
	Windows64           // windows-x64
	Windows32           // windows-x86
)

func (os *OS) MarshalText() ([]byte, error) {
	return []byte(os.String()), nil
}

func (os *OS) UnmarshalText(text []byte) error {
	s := string(text)
	for i := 0; i < len(_OS_index)-1; i++ {
		if OS(i).String() == s {
			*os = OS(i)
			return nil
		}
	}
	return errors.New("unknown os")
}
