//
// Copyright 2024 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package lsp

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/arduino/go-paths-helper"
	"go.bug.st/json"
)

type URI string

// DocumentURI Many of the interfaces contain fields that correspond to the URI of a document.
// For clarity, the type of such a field is declared as a DocumentUri. Over the wire, it will
// still be transferred as a string, but this guarantees that the contents of that string
// can be parsed as a valid URI.
type DocumentURI struct {
	url url.URL
}

// NilURI is the empty DocumentURI
var NilURI = DocumentURI{}

// for example, `"/c:"` or `"/A:"`
var expDriveWithLeadingSlashID = regexp.MustCompile("^/[a-zA-Z]:")

// for example, `"C:"` or `"A:"`
var expUppercaseDriveID = regexp.MustCompile("^[A-Z]:")

// AsPath convert the DocumentURI to a paths.Path
func (uri DocumentURI) AsPath() *paths.Path {
	return paths.New(uri.unbox()).Canonical()
}

// unbox convert the DocumentURI to a file path string
func (uri DocumentURI) unbox() string {
	path := uri.url.Path
	if expDriveWithLeadingSlashID.MatchString(path) {
		return path[1:]
	}
	return path
}

// Converts `"C:"` to `"c:"` to be compatible with VS Code URI's drive letter casing
// https://github.com/Microsoft/vscode/issues/68325#issuecomment-462239992
func lowercaseDriveSegment(pathSegment string) string {
	if expUppercaseDriveID.MatchString(pathSegment) {
		chars := []rune(pathSegment)
		chars[0] = unicode.ToLower(chars[0])
		return string(chars)
	}
	return pathSegment
}

func (uri DocumentURI) String() string {
	return uri.url.String()
}

// Ext returns the extension of the file pointed by the URI
func (uri DocumentURI) Ext() string {
	return filepath.Ext(uri.unbox())
}

// NewDocumentURIFromPath create a DocumentURI from the given Path object
func NewDocumentURIFromPath(path *paths.Path) DocumentURI {
	return NewDocumentURI(path.String())
}

var toSlash = filepath.ToSlash

// NewDocumentURI create a DocumentURI from the given string path
func NewDocumentURI(path string) DocumentURI {
	// transform path into URI
	path = toSlash(path)
	if len(path) == 0 || path[0] != '/' {
		path = "/" + path
	}
	segments := strings.Split(path, "/")
	encodedSegments := make([]string, len(segments))
	for i, segment := range segments {
		if len(segment) == 0 {
			encodedSegments[i] = segment
		} else {
			segment = lowercaseDriveSegment(segment)
			segment = url.QueryEscape(segment)
			// Spaces must be turned into `%20`. Otherwise, `url.QueryEscape`` encodes them to `+`.
			encodedSegments[i] = strings.ReplaceAll(segment, "+", "%20")
		}
	}
	urlPath := strings.Join(encodedSegments, "/")
	uri, err := NewDocumentURIFromURL("file://" + urlPath)
	if err != nil {
		panic(err)
	}
	return uri
}

// NewDocumentURIFromURL converts an URL into a DocumentURI
func NewDocumentURIFromURL(inURL string) (DocumentURI, error) {
	uri, err := url.Parse(inURL)
	if err != nil {
		return NilURI, err
	}
	return DocumentURI{url: *uri}, nil
}

// UnmarshalJSON implements json.Unmarshaller interface
func (uri *DocumentURI) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("expected JSON string for DocumentURI: %s", err)
	}

	newDocURI, err := NewDocumentURIFromURL(s)
	if err != nil {
		return fmt.Errorf("parsing DocumentURI: %s", err)
	}

	*uri = newDocURI
	return nil
}

func (uri *DocumentURI) UnmarshalText(text []byte) error {
	newDocURI, err := NewDocumentURIFromURL(string(text))
	if err != nil {
		return fmt.Errorf("parsing DocumentURI: %s", err)
	}

	*uri = newDocURI
	return nil
}

// MarshalJSON implements json.Marshaller interface
func (uri DocumentURI) MarshalJSON() ([]byte, error) {
	return json.Marshal(uri.url.String())
}

func (uri DocumentURI) MarshalText() (text []byte, err error) {
	return []byte(uri.String()), nil
}
