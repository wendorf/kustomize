package github

import (
	"testing"
)

func TestQueryFields(t *testing.T) {
	testCases := []struct {
		formatter queryField
		expected  string
	}{
		{
			formatter: Keyword("keyword"),
			expected:  "keyword",
		},
		{
			formatter: Filesize(RangeLessThan{23}),
			expected:  "size:<23",
		},
		{
			formatter: Filesize(RangeWithin{24, 64}),
			expected:  "size:24..64",
		},
		{
			formatter: Filesize(RangeGreaterThan{64}),
			expected:  "size:>64",
		},
		{
			formatter: Path("some/path/to/file"),
			expected:  "path:some/path/to/file",
		},
		{
			formatter: Filename("kustomization.yaml"),
			expected:  "filename:kustomization.yaml",
		},
	}

	for _, test := range testCases {
		if result := test.formatter.String(); result != test.expected {
			t.Errorf("got (%#v = %s), expected %s", test.formatter, result, test.expected)
		}
	}
}

func TestQueryType(t *testing.T) {
	testCases := []struct {
		query    Query
		expected string
	}{
		{
			query: QueryWith(
				Filesize(RangeWithin{24, 64}),
				Filename("kustomization.yaml"),
				Keyword("keyword1"),
				Keyword("keyword2"),
				Repo("user1/repo1"),
				User("user1"),
			),
			expected: "q=size:24..64+filename:kustomization.yaml+keyword1+keyword2+" +
				"repo:user1/repo1+user:user1",
		},
	}

	for _, test := range testCases {
		if queryStr := test.query.String(); queryStr != test.expected {
			t.Errorf("got (%#v = %s), expected %s", test.query, queryStr, test.expected)
		}

	}
}

func TestGithubSearchQuery(t *testing.T) {
	const (
		perPage = 100
	)

	testCases := []struct {
		rc                    RequestConfig
		codeQuery             Query
		fullRepoName          string
		path                  string
		expectedCodeQuery     string
		expectedContentsQuery string
		expectedCommitsQuery  string
	}{
		{
			rc: RequestConfig{
				perPage: perPage,
			},
			codeQuery: Query{
				Filename("kustomization.yaml"),
				Filesize(RangeWithin{64, 128}),
			},
			fullRepoName: "kubernetes-sigs/kustomize",
			path:         "examples/helloWorld/kustomization.yaml",

			expectedCodeQuery: "https://api.github.com/search/code?" +
				"q=filename:kustomization.yaml+size:64..128&order=desc&per_page=100&sort=indexed",

			expectedContentsQuery: "https://api.github.com/repos/kubernetes-sigs/kustomize/contents/" +
				"examples/helloWorld/kustomization.yaml?per_page=100",

			expectedCommitsQuery: "https://api.github.com/repos/kubernetes-sigs/kustomize/commits?" +
				"path=examples%2FhelloWorld%2Fkustomization.yaml&per_page=100",
		},
		{
			rc: RequestConfig{
				perPage: perPage,
			},
			codeQuery: Query{
				Filename("kustomization.yaml"),
				Filesize(RangeWithin{64, 128}),
			},
			fullRepoName: "kubernetes-sigs/kustomize",
			path:         "examples 1/helloWorld/kustomization.yaml",

			expectedCodeQuery: "https://api.github.com/search/code?" +
				"q=filename:kustomization.yaml+size:64..128&order=desc&per_page=100&sort=indexed",

			expectedContentsQuery: "https://api.github.com/repos/kubernetes-sigs/kustomize/contents/" +
				"examples%201/helloWorld/kustomization.yaml?per_page=100",

			expectedCommitsQuery: "https://api.github.com/repos/kubernetes-sigs/kustomize/commits?" +
				"path=examples+1%2FhelloWorld%2Fkustomization.yaml&per_page=100",
		},
	}

	for _, test := range testCases {
		if result := test.rc.CodeSearchRequestWith(test.codeQuery).URL(); result != test.expectedCodeQuery {
			t.Errorf("Got code query: %s, expected %s", result, test.expectedCodeQuery)
		}

		if result := test.rc.ContentsRequest(test.fullRepoName, test.path); result != test.expectedContentsQuery {
			t.Errorf("Got contents query: %s, expected %s", result, test.expectedContentsQuery)
		}
		if result := test.rc.CommitsRequest(test.fullRepoName, test.path); result != test.expectedCommitsQuery {
			t.Errorf("Got commits query: %s, expected %s", result, test.expectedCommitsQuery)
		}
	}
}
