package atom

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"testing"
)

func TestMarshal(t *testing.T) {
	feed, err := testLoadAtomFile("atom_1.0_all.xml")
	if err != nil {
		t.Fatal(err)
		return
	}

	valErr := feed.Validate()
	if valErr != nil {
		t.Fatal(err)
		return
	}

	original, err := ioutil.ReadFile("samples/atom_1.0_all.xml")
	if err != nil {
		t.Fatal(err)
		return
	}

	data, err := xml.MarshalIndent(feed, "", "    ")
	if err != nil {
		t.Fatal(err)
		return
	}

	if bytes.Compare(data, original) != 0 {
		t.Fatal("Marshaled output differens from input")
	}
}

func TestAllFeedElements(t *testing.T) {
	feed, err := testLoadAtomFile("atom_1.0_all.xml")
	if err != nil {
		t.Fatal(err)
		return
	}

	if len(feed.Author) != 2 ||
		feed.Author[0].Name != "John Doe" ||
		feed.Author[0].URI != "http://john.doe" ||
		feed.Author[0].Email != "john@test.com" ||
		feed.Author[1].Name != "Jane Doe" ||
		feed.Author[1].URI != "http://jane.doe" ||
		feed.Author[1].Email != "jane@test.com" {
		t.Fatal("Author elements did not parse successfully.")
	}

	if len(feed.Category) != 2 ||
		feed.Category[0].Term != "testcat" ||
		feed.Category[0].Scheme != "http://testcat.john.doe" ||
		feed.Category[0].Label != "Test Category" ||
		feed.Category[1].Term != "testcat2" ||
		feed.Category[1].Scheme != "http://testcat2.john.doe" ||
		feed.Category[1].Label != "Test Category 2" {
		t.Fatal("Category elements did not parse successfully.")
	}

	if len(feed.Contributor) != 2 ||
		feed.Contributor[0].Name != "Contrib1" ||
		feed.Contributor[0].URI != "http://contrib1" ||
		feed.Contributor[0].Email != "contrib1@test.com" ||
		feed.Contributor[1].Name != "Contrib2" ||
		feed.Contributor[1].URI != "http://contrib2" ||
		feed.Contributor[1].Email != "contrib2@test.com" {
		t.Fatal("Contributor elements did not parse successfully.")
	}

	if feed.Generator == nil ||
		feed.Generator.URI != "http://github.com/test/atom" ||
		feed.Generator.Version != "1.0" ||
		feed.Generator.Text != "Golang ATOM package" {
		t.Fatal("Generator element did not parse successfully.")
	}

	if feed.Icon != "http://icon.test.com" {
		t.Fatal("Icon element did not parse successfully.")
	}

	if feed.ID != "http://www.test.com/blog" {
		t.Fatal("ID element did not parse successfully.")
	}

	if len(feed.Link) != 2 ||
		feed.Link[0].Rel != "alternate" ||
		feed.Link[0].Href != "http://www.test.com/blog2" ||
		feed.Link[0].HrefLang != "de" ||
		feed.Link[0].Type != "application/xml" ||
		feed.Link[0].Title != "German feed" ||
		feed.Link[0].Length != "10" ||
		feed.Link[1].Rel != "self" ||
		feed.Link[1].Href != "http://www.test.com/blog" {
		t.Fatal("Link element did not parse successfully.")
	}

	if feed.Logo != "http://www.test.com/logo" {
		t.Fatal("Logo element did not parse successfully.")
	}

	if feed.Rights != "Test Corp TM" {
		t.Fatal("Rights element did not parse successfully.")
	}

	if feed.Subtitle != "Test subtitle" {
		t.Fatal("Subtitle element did not parse successfully.")
	}

	if feed.Title != "Test feed" {
		t.Fatal("Title element did not parse successfully.")
	}

	if feed.Updated != "2006-11-04T09:11:03-08:00" {
		t.Fatal("Updated element did not parse successfully.")
	}
}
