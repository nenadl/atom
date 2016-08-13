package atom

import (
	"encoding/xml"
	"time"
)

// Common ATOM xml attributes. Validation:
//   - Base if not empty must be a valid url.URL.
type Common struct {
	Base string `xml:"http://www.w3.org/XML/1998/namespace base,attr,omitempty"`
	Lang string `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
}

// Text is entry content. Validation:
//   - Type if present must be one of "text", "html", or "xhtml". If this text
//     is however Entry.Content, then any valid MIME type is also acceptable.
//   - Src if present must be a valid url.URL.
type Text struct {
	Common

	Type string `xml:"type,attr,omitempty"`
	Src  string `xml:"src,attr,omitempty"`
	Body string `xml:",innerxml"`
}

// Extension represents a custom atom element.
type Extension struct {
	XMLName xml.Name
	XML     string `xml:",innerxml"`
}

// Person is a feed author or contributor. I.e. <author name="John Doe"
// uri="http://www.johndoe.com" email="john@test.com" />.
//
// Validation:
//   - URI if present must be a valid url.URL.
//   - Email if present must be a valid email address. mail.ParseAddress is
//     used to validate the email.
type Person struct {
	Common

	Name  string `xml:"name"`
	URI   string `xml:"uri,omitempty"`
	Email string `xml:"email,omitempty"`

	Extension []Extension `xml:",any,omitempty"`
}

// TimeStr is a date/time attribute that is in the correct
// https://www.ietf.org/rfc/rfc3339.txt format.
type TimeStr string

// AtomTimeFormat is the https://www.ietf.org/rfc/rfc3339.txt format
const AtomTimeFormat = time.RFC3339

// Time creates a TimeStr from a time.Time object with the correct formatting.
func Time(t time.Time) TimeStr {
	return TimeStr(t.Format(AtomTimeFormat))
}

// TimeParse parses an atom time to a golang time.Time value
func TimeParse(t TimeStr) (time.Time, error) {
	return time.Parse(AtomTimeFormat, string(t))
}

// Category is the category that this feed or entry belongs to. Validation:
//   - Term must be present and can't be empty.
//   - Scheme if present must be a url.URL.
type Category struct {
	Common

	Term   string `xml:"term,attr"`
	Scheme string `xml:"scheme,attr,omitempty"`
	Label  string `xml:"label,attr,omitempty"`
}

// Generator is the generating agent for this feed. Validation:
//   - if URI is present it must be a valid url.URL.
type Generator struct {
	Common

	URI     string `xml:"uri,attr"`
	Version string `xml:"version,attr"`
	Text    string `xml:",chardata"`
}

// Link is a feed or entry link, i.e: <link rel="enclosure" type="audio/mpeg"
// title="MP3" href="http://www.example.org/myaudiofile.mp3" hreflang="de"
// length="1234" />. Validation:
//   - Href must be a valid url.URL.
type Link struct {
	Common

	Href     string `xml:"href,attr"`
	Rel      string `xml:"rel,attr,omitempty"`
	Type     string `xml:"type,attr,omitempty"`
	HrefLang string `xml:"hreflang,attr,omitempty"`
	Title    string `xml:"title,attr,omitempty"`
	Length   string `xml:"length,attr,omitempty"`
}

// Source is the verbatim source feed that this feed was copied from. Note that
// the embedded feed should not contain any entries. Only the feed metadata is
// embedded. Validations:
//   - Feed.Entry must be empty
type Source struct {
	Feed *Feed
}

// Entry is a single entry inside a Feed. Validation:
//   - ID can't be empty.
//   - Can't contain more than one link element of rel="alternate" and the same values.
//   - Title can't be empty.
//   - Updated can't be empty and must be a valid time string.
type Entry struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom entry"`
	Common

	// Metadata
	Author      []Person   `xml:"author"`
	Category    []Category `xml:"category,omitempty"`
	Content     *Text      `xml:"content,omitempty"`
	Contributor []Person   `xml:"contributor,omitempty"`
	ID          string     `xml:"id"`
	Link        []Link     `xml:"link,omitempty"`
	Published   TimeStr    `xml:"published,omitempty"`
	Rights      string     `xml:"rights,omitempty"`
	Source      *Source    `xml:"source,omitempty"`
	Summary     *Text      `xml:"summary,omitempty"`
	Title       string     `xml:"title"`
	Updated     TimeStr    `xml:"updated"`

	// Custom elements
	Extension []Extension `xml:",any,omitempty"`
}

// Feed is the top level ATOM syndication element. Validation:
//   - Author can't be empty, unless present in all entry elements.
//   - Can't contain more than one link element of rel="alternate" and the same values.
//   - Title can't be empty.
//   - Icon if present must be a valid url.URL.
//   - ID can't be empty.
//   - Logo if present must be a valid url.URL.
//   - Updated can't be empty and must be a valid time string.
type Feed struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Common

	// Metadata
	Author      []Person   `xml:"author"`
	Category    []Category `xml:"category,omitempty"`
	Contributor []Person   `xml:"contributor,omitempty"`
	Generator   *Generator `xml:"generator,omitempty"`
	Icon        string     `xml:"icon,omitempty"`
	ID          string     `xml:"id"`
	Link        []Link     `xml:"link,omitempty"`
	Logo        string     `xml:"logo,omitempty"`
	Rights      string     `xml:"rights,omitempty"`
	Subtitle    string     `xml:"subtitle,omitempty"`
	Title       string     `xml:"title"`
	Updated     TimeStr    `xml:"updated"`

	// Entries
	Entry []Entry `xml:"entry"`

	// Custom elements
	Extension []Extension `xml:",any,omitempty"`
}
