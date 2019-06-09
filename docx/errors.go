package docx

// WrongXMLSlice returned if not valid XmlSlice is passed
type WrongXMLSlice struct {
	XMLSlice string
	Reason   string
}

func (e *WrongXMLSlice) Error() string {
	return e.Reason
}
