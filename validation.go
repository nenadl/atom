package atom

import (
	"net/mail"
	"net/url"
	"time"
)

// ValidationError is an error encountered while trying to validate an
// ATOM feed.
type ValidationError struct {
	issues []string
}

func (err *ValidationError) Error() string {
	result := "An ATOM validation error occured:\n"

	for _, issue := range err.issues {
		result += "- " + issue + "\n"
	}

	return result
}

// Merge validate errors. other must be a *ValidationError.
func (err *ValidationError) Merge(other error) {
	if other == nil {
		return
	}

	err.issues = append(err.issues, other.(*ValidationError).issues...)
}

// Add Adds a string to the list of issues
func (err *ValidationError) Add(issue string) {
	err.issues = append(err.issues, issue)
}

// AddErr Adds an error to the list of validaiton issues. It converts it to a
// string first.
func (err *ValidationError) AddErr(issue error) {
	err.issues = append(err.issues, issue.Error())
}

// NilIfEmpty returns nil if there are no issues in the ValidationError,
// otherwise it returns itself.
func (err *ValidationError) NilIfEmpty() error {
	if len(err.issues) == 0 {
		return nil
	}

	return err
}

// Validate the common ATOM structure.
func validateCommon(common Common) error {
	var errs = new(ValidationError)

	if common.Base != "" {
		if _, err := url.Parse(common.Base); err != nil {
			errs.AddErr(err)
		}
	}

	return errs.NilIfEmpty()
}

// Validate the text ATOM structure.
func (text *Text) validate(isContent bool) error {
	var errs = new(ValidationError)

	errs.Merge(validateCommon(text.Common))

	if !isContent {
		if text.Type != "" {
			if text.Type != "text" && text.Type != "html" && text.Type != "xhtml" {
				errs.Add("Text.Type must be either: text, html or xhtml.")
			}
		}
	}

	if text.Src != "" {
		if _, err := url.Parse(text.Src); err != nil {
			errs.AddErr(err)
		}
	}

	return errs.NilIfEmpty()
}

// Validate the person ATOM structure.
func (person *Person) validate() error {
	var errs = new(ValidationError)

	errs.Merge(validateCommon(person.Common))

	if person.URI != "" {
		if _, err := url.Parse(person.URI); err != nil {
			errs.AddErr(err)
		}
	}

	if person.Email != "" {
		if _, err := mail.ParseAddress(person.Email); err != nil {
			errs.AddErr(err)
		}
	}

	return errs.NilIfEmpty()
}

// Validate the time format.
func (timeStr *TimeStr) validate() error {
	var errs = new(ValidationError)

	if _, err := time.Parse(AtomTimeFormat, string(*timeStr)); err != nil {
		errs.AddErr(err)
	}

	return errs.NilIfEmpty()
}

// Validate an ATOM category.
func (category *Category) validate() error {
	var errs = new(ValidationError)

	errs.Merge(validateCommon(category.Common))

	if category.Scheme != "" {
		if _, err := url.Parse(category.Scheme); err != nil {
			errs.AddErr(err)
		}
	}

	return errs.NilIfEmpty()
}

// Validate an ATOM generator.
func (generator *Generator) validate() error {
	var errs = new(ValidationError)

	errs.Merge(validateCommon(generator.Common))

	if generator.URI != "" {
		if _, err := url.Parse(generator.URI); err != nil {
			errs.AddErr(err)
		}
	}

	return errs.NilIfEmpty()
}

// Validate an ATOM link.
func (link *Link) validate() error {
	var errs = new(ValidationError)

	errs.Merge(validateCommon(link.Common))

	if _, err := url.Parse(link.Href); err != nil {
		errs.AddErr(err)
	}

	return errs.NilIfEmpty()
}

// Validate validates an entry document
func (entry *Entry) Validate() error {
	var errs = new(ValidationError)

	errs.Merge(validateCommon(entry.Common))

	for _, author := range entry.Author {
		errs.Merge(author.validate())
	}

	for _, cat := range entry.Category {
		errs.Merge(cat.validate())
	}

	if entry.Content != nil {
		errs.Merge(entry.Content.validate(true))
	}

	for _, contrib := range entry.Contributor {
		errs.Merge(contrib.validate())
	}

	if entry.ID == "" {
		errs.Add("Entry.ID can't be empty.")
	}

	foundAlternate := false
	for _, link := range entry.Link {
		errs.Merge(link.validate())

		if link.Rel == "alternate" {
			if foundAlternate {
				errs.Add("Only one Feed.Link with rel=\"alternate\" can exist.")
			} else {
				foundAlternate = true
			}
		}
	}

	if entry.Published != "" {
		errs.Merge(entry.Published.validate())
	}

	if entry.Source != nil {
		if len(entry.Source.Feed.Entry) > 0 {
			errs.Add("Entry.Source can't contain any entries.")
		}
	}

	if entry.Summary != nil {
		errs.Merge(entry.Summary.validate(false))
	}

	if entry.Title == "" {
		errs.Add("Entry.Title can't be empty.")
	}

	if entry.Updated != "" {
		errs.Merge(entry.Updated.validate())
	} else {
		errs.Add("Entry.Updated can't be empty.")
	}

	if entry.Source != nil {
		errs.Merge(entry.Source.Feed.Validate())
	}

	return errs.NilIfEmpty()
}

// Validate check the atom feed structure to make sure it's comformant with the
// ATOM spec: https://tools.ietf.org/html/rfc4287. All validation errors are
// returned together.
func (feed *Feed) Validate() error {
	var errs = new(ValidationError)

	errs.Merge(validateCommon(feed.Common))

	for _, author := range feed.Author {
		errs.Merge(author.validate())
	}

	// Must be in all entries if not in feed.
	if len(feed.Author) == 0 {
		for _, entry := range feed.Entry {
			if len(entry.Author) == 0 {
				errs.Add("Author must be present in Feed or in every Entry.")
				break
			}
		}
	}

	for _, cat := range feed.Category {
		errs.Merge(cat.validate())
	}

	for _, contrib := range feed.Contributor {
		errs.Merge(contrib.validate())
	}

	if feed.Generator != nil {
		errs.Merge(feed.Generator.validate())
	}

	if feed.Icon != "" {
		if _, err := url.Parse(feed.Icon); err != nil {
			errs.AddErr(err)
		}
	}

	if feed.ID == "" {
		errs.Add("Feed.ID can't be empty.")
	}

	foundAlternate := false
	for _, link := range feed.Link {
		errs.Merge(link.validate())

		if link.Rel == "alternate" {
			if foundAlternate {
				errs.Add("Only one Feed.Link with rel=\"alternate\" can exist.")
			} else {
				foundAlternate = true
			}
		}
	}

	if feed.Logo != "" {
		if _, err := url.Parse(feed.Logo); err != nil {
			errs.AddErr(err)
		}
	}

	if feed.Title == "" {
		errs.Add("Feed.Title can't be empty.")
	}

	if feed.Updated == "" {
		errs.Add("Feed.Updated can't be empty.")
	}

	if feed.Updated != "" {
		errs.Merge(feed.Updated.validate())
	} else {
		errs.Add("Feed.Updated can't be empty.")
	}

	for _, entry := range feed.Entry {
		errs.Merge(entry.Validate())
	}

	return errs.NilIfEmpty()
}
