package gitlab

import "testing"

func TestProjects_ByPath(t *testing.T) {
	before(t)

	proj := createProject(t, "temporary-awesome-", "Temporary repository to check if exists")
	defer deleteProject(t, proj)

	type _test struct {
		path string
	}
	tests := []*_test{
		&_test{*proj.PathWithNamespace},
	}
	for _, test := range tests {
		proj, err := GitLabClient.Projects.ByPath(test.path)
		if err != nil {
			t.Fatal(err)
		}
		if *proj.PathWithNamespace != test.path {
			t.Errorf("expecting '%s', got '%s'\n", test.path, *proj.PathWithNamespace)
		}
	}
	tests = []*_test{
		&_test{"root/nonexistingrepo"},
	}
	for _, test := range tests {
		proj, err := GitLabClient.Projects.ByPath(test.path)
		if _, ok := err.(*NotFound); !ok {
			t.Fatalf("expecting not found, got: %s", *proj.PathWithNamespace)
		}
	}
}
