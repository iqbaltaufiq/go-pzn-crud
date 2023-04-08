package helper

// error checking
// to shorten code from 5 lines to one
func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
