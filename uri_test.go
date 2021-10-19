package lsp

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUriToPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		require.Equal(t, "C:\\Users\\test\\Sketch.ino", DocumentURI("file:///C:/Users/test/Sketch.ino").Unbox())
		require.Equal(t, "C:\\Users\\test\\Sketch.ino", DocumentURI("file:///c%3A/Users/test/Sketch.ino").Unbox())
	} else {
		require.Equal(t, "/Users/test/Sketch.ino", DocumentURI("file:///Users/test/Sketch.ino").Unbox())
	}
	require.Equal(t, string(filepath.Separator)+"\U0001F61B", DocumentURI("file:///%25F0%259F%2598%259B").Unbox())
}

func TestPathToUri(t *testing.T) {
	if runtime.GOOS == "windows" {
		require.Equal(t, DocumentURI("file:///C:/Users/test/Sketch.ino"), NewDocumentURI("C:\\Users\\test\\Sketch.ino"))
	} else {
		require.Equal(t, DocumentURI("file:///Users/test/Sketch.ino"), NewDocumentURI("/Users/test/Sketch.ino"))
	}
	require.Equal(t, DocumentURI("file:///%25F0%259F%2598%259B"), NewDocumentURI("\U0001F61B"))
}
