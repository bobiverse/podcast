package podcast

// Owner ..
type Owner struct {
	Text  string `xml:",chardata"`
	Name  string `xml:"name"`
	Email string `xml:"email"`
}
