package tinyfs

import (
	"fmt"
	"testing"
)

func TestAbsPath(t *testing.T) {
	data := []struct {
		pwd           string
		path          string
		leadingSlash  bool
		trailingSlash bool
		wantPath      string
	}{
		{},
		{
			leadingSlash: true,
			wantPath:     "/",
		},
		{
			trailingSlash: true,
			wantPath:      "/",
		},
		{
			leadingSlash:  true,
			trailingSlash: true,
			wantPath:      "/",
		},
		{
			path: ".",
		},
		{
			path:         ".",
			leadingSlash: true,
			wantPath:     "/",
		},
		{
			path:          ".",
			trailingSlash: true,
			wantPath:      "/",
		},
		{
			path:          ".",
			leadingSlash:  true,
			trailingSlash: true,
			wantPath:      "/",
		},
		{
			path: "/",
		},
		{
			path:         "/",
			leadingSlash: true,
			wantPath:     "/",
		},
		{
			path:          "/",
			trailingSlash: true,
			wantPath:      "/",
		},
		{
			path:          "/",
			leadingSlash:  true,
			trailingSlash: true,
			wantPath:      "/",
		},
		{
			path: "..",
		},
		{
			path:         "..",
			leadingSlash: true,
			wantPath:     "/",
		},
		{
			path:          "..",
			trailingSlash: true,
			wantPath:      "/",
		},
		{
			path:          "..",
			leadingSlash:  true,
			trailingSlash: true,
			wantPath:      "/",
		},
		{
			path:     "foo",
			wantPath: "foo",
		},
		{
			path:         "foo",
			leadingSlash: true,
			wantPath:     "/foo",
		},
		{
			path:          "foo",
			trailingSlash: true,
			wantPath:      "foo/",
		},
		{
			path:          "foo",
			leadingSlash:  true,
			trailingSlash: true,
			wantPath:      "/foo/",
		},
		{
			pwd:      "foo/bar",
			path:     "baz/buz",
			wantPath: "foo/bar/baz/buz",
		},
		{
			pwd:          "foo/bar",
			path:         "baz/buz",
			leadingSlash: true,
			wantPath:     "/foo/bar/baz/buz",
		},
		{
			pwd:           "foo/bar",
			path:          "baz/buz",
			trailingSlash: true,
			wantPath:      "foo/bar/baz/buz/",
		},
		{
			pwd:           "foo/bar",
			path:          "baz/buz",
			leadingSlash:  true,
			trailingSlash: true,
			wantPath:      "/foo/bar/baz/buz/",
		},
	}

	for _, d := range data {
		t.Run(fmt.Sprintf("%+v", d), func(t *testing.T) {
			got := absPath(d.pwd, d.path, d.leadingSlash, d.trailingSlash)
			if got != d.wantPath {
				t.Errorf("got %v, want %v", got, d.wantPath)
			}
		})
	}
}
