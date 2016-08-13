package atom

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	"encoding/xml"

	"golang.org/x/net/html/charset"
)

var (
	testFixtures = []string{
		// All features
		"atom_1.0_all.xml",

		// From Java Rome
		"atom_1.0_b.xml",
		"atom_1.0_bray.xml",
		"atom_1.0_prefix.xml",
		"atom_1.0_ruby.xml",

		// From http://docs.oasis-open.org/cmis/CMIS/v1.1/os/examples/atompub/
		"doQuery-response.xml",
		"getAllVersions-response.xml",
		"getChildren-response.xml",
		"getDescendants-response.xml",
		"getObjectParents-response.xml",
		"getTypeChildren-response.xml",

		// Slashdot
		"slashdot.xml",
	}
)

func testLoadAtomFile(fileName string) (*Feed, error) {
	data, err := ioutil.ReadFile("samples/" + fileName)
	if err != nil {
		return nil, errors.New(fileName + ": " + err.Error())
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel

	var feed Feed
	err = decoder.Decode(&feed)
	if err != nil {
		return nil, errors.New(fileName + ": " + err.Error())
	}

	return &feed, nil
}

func TestWithValidAtom(t *testing.T) {
	for _, fileName := range testFixtures {
		feed, err := testLoadAtomFile(fileName)
		if err != nil {
			t.Error("Load file error (" + fileName + "): " + err.Error())
			continue
		}

		err = feed.Validate()
		if err != nil {
			t.Error("Validation error (" + fileName + "): " + err.Error())
		}
	}
}

func testLoadAtomFileAndCheckError(t *testing.T, fileName string, errToCheck string) {
	feed, err := testLoadAtomFile(fileName)
	if err != nil {
		t.Error("Load file error (" + fileName + "): " + err.Error())
		return
	}

	validationError := feed.Validate()
	if validationError == nil {
		t.Error("Validation failed on bad feed - it was successful, file: " + fileName)
		return
	}

	if !strings.Contains(validationError.Error(), errToCheck) {
		t.Error("Validation failed on bad feed: " + fileName +
			" , error: " + errToCheck +
			" result: " + validationError.Error())
	}
}

func TestFeedValidations(t *testing.T) {
	testLoadAtomFileAndCheckError(t, "bad_atom_author.xml", "Author must be present in Feed or in every Entry.")
	testLoadAtomFileAndCheckError(t, "bad_atom_id.xml", "Feed.ID can't be empty")
	testLoadAtomFileAndCheckError(t, "bad_atom_link.xml", "Only one Feed.Link with rel=\"alternate\" can exist.")
	testLoadAtomFileAndCheckError(t, "bad_atom_title.xml", "Feed.Title can't be empty.")
	testLoadAtomFileAndCheckError(t, "bad_atom_updated.xml", "Feed.Updated can't be empty.")
}

func TestOtherValidations(t *testing.T) {
	testLoadAtomFileAndCheckError(t, "bad_atom_text.xml", "Text.Type must be either: text, html or xhtml.")
	testLoadAtomFileAndCheckError(t, "bad_atom_url.xml", "invalid URL escape")
	testLoadAtomFileAndCheckError(t, "bad_atom_time.xml", "cannot parse")
	testLoadAtomFileAndCheckError(t, "bad_atom_email.xml", "mail: missing phrase")
}
