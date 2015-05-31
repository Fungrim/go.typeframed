package typeframed

import ()

type HeaderCapture interface {
	Capture(bytes []byte) error
}
