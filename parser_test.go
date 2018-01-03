package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ParseSpec struct {
	Input  string
	Output []Mapping
	Error  error
}

type ParseLineSpec struct {
	Title  string
	Input  string
	Output Mapping
	Ok     bool
}

func TestParse(t *testing.T) {
	specs := []ParseSpec{
		ParseSpec{
			Input: `* @user1

# Some comment
foo/* @user2 email@example.com

### More comments ###

   bar/qux/*  	  @user3     @org/team1 # @notme
`,
			Output: []Mapping{
				Mapping{
					Path: "*",
					Owners: []Owner{
						Owner("@user1"),
					},
				},
				Mapping{
					Path: "foo/*",
					Owners: []Owner{
						Owner("@user2"),
						Owner("email@example.com"),
					},
				},
				Mapping{
					Path: "bar/qux/*",
					Owners: []Owner{
						Owner("@user3"),
						Owner("@org/team1"),
					},
				},
			},
			Error: nil,
		},
	}
	for _, s := range specs {
		output, err := Parse(strings.NewReader(s.Input))
		assert.Equal(t, s.Error, err)
		assert.Equal(t, s.Output, output)
	}
}

func TestParseLine(t *testing.T) {
	specs := []ParseLineSpec{
		ParseLineSpec{
			Title: "maps a single Github User",
			Input: `* @nathanleiby`,
			Output: Mapping{
				Path: "*",
				Owners: []Owner{
					Owner("@nathanleiby"),
				},
			},
			Ok: true,
		},
		ParseLineSpec{
			Title: "allows email addresses",
			Input: `* nathan.leiby@clever.com`,
			Output: Mapping{
				Path: "*",
				Owners: []Owner{
					Owner("nathan.leiby@clever.com"),
				},
			},
			Ok: true,
		},
		ParseLineSpec{
			Title: "maps multiple Users",
			Input: `* @nathanleiby @clever/some-team nathan.leiby@clever.com`,
			Output: Mapping{
				Path: "*",
				Owners: []Owner{
					Owner("@nathanleiby"),
					Owner("@clever/some-team"),
					Owner("nathan.leiby@clever.com"),
				},
			},
			Ok: true,
		},

		ParseLineSpec{
			Title: "ignores invalid users",
			Input: `* @nathanleiby @user2 /// 12412 aa @user3`,
			Output: Mapping{
				Path: "*",
				Owners: []Owner{
					Owner("@nathanleiby"),
					Owner("@user2"),
					Owner("@user3"),
				},
			},
			Ok: true,
		},
		ParseLineSpec{
			Title: "stops reading line once there's a comment symbol (#)",
			Input: `*        @nathanleiby        # @ignoreduser`,
			Output: Mapping{
				Path: "*",
				Owners: []Owner{
					Owner("@nathanleiby"),
				},
			},
			Ok: true,
		},
		ParseLineSpec{
			Title: "ignores irrelevant whitespace",
			Input: "  \t   *\t	@nathanleiby  \t   @user2    \t\t #\t  @user3",
			Output: Mapping{
				Path: "*",
				Owners: []Owner{
					Owner("@nathanleiby"),
					Owner("@user2"),
				},
			},
			Ok: true,
		},
		ParseLineSpec{
			Title: "handles other paths",
			Input: "foo/*	@nathanleiby",
			Output: Mapping{
				Path: "foo/*",
				Owners: []Owner{
					Owner("@nathanleiby"),
				},
			},
			Ok: true,
		},
		ParseLineSpec{
			Title:  "ignores lines with comments",
			Input:  `# somecomment`,
			Output: Mapping{},
			Ok:     false,
		},
	}
	for _, s := range specs {
		owner, ok := parseLine(s.Input)
		assert.Equal(t, s.Ok, ok)
		assert.Equal(t, s.Output, owner)
	}
}
