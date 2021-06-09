package utils

//ErrTranscribe indicates transcription error
type ErrTranscribe struct {
	Msg string
}

//NewErrTranscribe creates new error
func NewErrTranscribe(msg string) *ErrTranscribe {
	return &ErrTranscribe{Msg: msg}
}

func (e *ErrTranscribe) Error() string {
	return e.Msg
}
