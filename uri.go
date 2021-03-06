//
// Copyright 2021 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package lsp

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"

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

var expDriveID = regexp.MustCompile("^/[a-zA-Z]:")

// AsPath convert the DocumentURI to a paths.Path
func (uri DocumentURI) AsPath() *paths.Path {
	return paths.New(uri.unbox()).Canonical()
}

// unbox convert the DocumentURI to a file path string
func (uri DocumentURI) unbox() string {
	path := uri.url.Path
	if expDriveID.MatchString(path) {
		return path[1:]
	}
	return path
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
	// tranform path into URI
	path = toSlash(path)
	if len(path) == 0 || path[0] != '/' {
		path = "/" + path
	}
	uri, err := NewDocumentURIFromURL("file://" + path)
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
		return fmt.Errorf("expoected JSON string for DocumentURI: %s", err)
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
