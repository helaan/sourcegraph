package errcode

// A PresentationError is an error with a message (returned by the PresentationError method) that is
// suitable for presentation to the user.
type PresentationError interface {
	error

	// PresentationError returns the message suitable for presentation to the user. The message
	// should be written in full sentences and must not contain any information that the user is not
	// authorized to see.
	PresentationError() string
}

// WithPresentationMessage annotates err with a new message suitable for presentation to the
// user. If err is nil, WithPresentationMessage returns nil. Otherwise, the return value implements
// PresentationError.
//
// The message should be written in full sentences and must not contain any information that the
// user is not authorized to see.
func WithPresentationMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &presentationError{cause: err, msg: message}
}

// NewPresentationError returns a new error with a message suitable for presentation to the user.
// The message should be written in full sentences and must not contain any information that the
// user is not authorized to see.
//
// If there is an underlying error associated with this message, use WithPresentationMessage
// instead.
func NewPresentationError(message string) error {
	return &presentationError{cause: nil, msg: message}
}

// presentationError implements PresentationError.
type presentationError struct {
	cause error
	msg   string
}

func (e *presentationError) Error() string {
	if e.cause != nil {
		return e.cause.Error()
	}
	return e.msg
}

func (e *presentationError) PresentationError() string { return e.msg }

// PresentationMessage returns the message, if any, suitable for presentation to the user for err or
// one of its causes. An error provides a presentation message by implementing the PresentationError
// interface (e.g., by using WithPresentationMessage). If no presentation message exists for err,
// the empty string is returned.
func PresentationMessage(err error) string {
	type causer interface {
		Cause() error
	}

	for err != nil {
		if e, ok := err.(PresentationError); ok {
			return e.PresentationError()
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return ""
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_775(size int) error {
	const bufSize = 1024

	f, err := os.Create("/tmp/test")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fb := bufio.NewWriter(f)
	defer fb.Flush()

	buf := make([]byte, bufSize)

	for i := size; i > 0; i -= bufSize {
		if _, err = rand.Read(buf); err != nil {
			fmt.Printf("error occurred during random: %!s(MISSING)\n", err)
			break
		}
		bR := bytes.NewReader(buf)
		if _, err = io.Copy(fb, bR); err != nil {
			fmt.Printf("failed during copy: %!s(MISSING)\n", err)
			break
		}
	}

	return err
}		
