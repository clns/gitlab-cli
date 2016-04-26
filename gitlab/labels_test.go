package gitlab

import (
	"testing"

	"path"
	"runtime"

	"reflect"

	gogitlab "github.com/xanzy/go-gitlab"
)

func TestLabels_UpdateWithRegex(t *testing.T) {
	before(t)

	proj := createProject(t, "temporary-update-labels-", "Temporary repository to update labels into")
	defer deleteProject(t, proj)

	addLabel(t, proj, "category#label", "#000000")
	addLabel(t, proj, "misc#anotherlabel", "#ff0000")

	// update with regex
	if err := GitLabClient.Labels.UpdateWithRegex(*proj.ID, &gogitlab.UpdateLabelOptions{
		Name:    "(.+)#(.+)",
		NewName: "${1}/${2}",
	}); err != nil {
		t.Fatal(err)
	}

	labelsExist(t, proj, []*gogitlab.Label{
		&gogitlab.Label{"category/label", ""},
		&gogitlab.Label{"misc/anotherlabel", ""},
	})

	// update again without regex
	if err := GitLabClient.Labels.UpdateWithRegex(*proj.ID, &gogitlab.UpdateLabelOptions{
		Name:    "category/label",
		NewName: "category-label",
	}); err != nil {
		t.Fatal(err)
	}

	labelsExist(t, proj, []*gogitlab.Label{
		&gogitlab.Label{"category-label", ""},
		&gogitlab.Label{"misc/anotherlabel", ""},
	})

	// update color
	if err := GitLabClient.Labels.UpdateWithRegex(*proj.ID, &gogitlab.UpdateLabelOptions{
		Name:  "^misc",
		Color: "#ff7863",
	}); err != nil {
		t.Fatal(err)
	}

	labelsExist(t, proj, []*gogitlab.Label{
		&gogitlab.Label{"misc/anotherlabel", "#ff7863"},
	})
}

func TestLabels_DeleteWithRegex(t *testing.T) {
	before(t)

	proj := createProject(t, "temporary-delete-labels-from-", "Temporary repository to delete labels from")
	defer deleteProject(t, proj)

	addLabel(t, proj, "test-label", "#000000")

	if err := GitLabClient.Labels.DeleteWithRegex(*proj.ID, ""); err != nil {
		t.Fatal(err)
	}
	if labels := getLabels(t, *proj.ID); len(labels) > 0 {
		t.Fatalf("labels still exist after supposedly deleting them all: %v", labels)
	}
}

func TestLabels_CopyGlobalLabelsTo(t *testing.T) {
	before(t)

	proj := createProject(t, "temporary-copy-globals-to-", "Temporary repository to copy global labels to")
	defer deleteProject(t, proj)

	globalLabels := getLabels(t, *proj.ID)

	if err := GitLabClient.Labels.DeleteWithRegex(*proj.ID, ""); err != nil {
		t.Fatal(err)
	}

	if err := GitLabClient.Labels.CopyGlobalLabelsTo(*proj.ID); err != nil {
		t.Fatal(err)
	}

	labels := getLabels(t, *proj.ID)
	if len(labels) != len(globalLabels) {
		t.Fatalf("different number of labels\nglobalLabels: %v\nrepoLabels: %v", globalLabels, labels)
	}
	for i, label := range labels {
		global := globalLabels[i]
		if label.Name != global.Name || label.Color != global.Color {
			t.Fatalf("labels are different\nglobalLabels: %v\nrepoLabels: %v", globalLabels, labels)
		}
	}
}

// Helper functions:

func getLabels(tb testing.TB, pid interface{}) []*gogitlab.Label {
	labels, _, err := GitLabClient.Labels.ListLabels(pid)
	if err != nil {
		// The failure happens at wherever we were called, not here
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			tb.Fatalf("Unable to get caller")
		}
		tb.Fatalf("%s:%v %v", path.Base(file), line, err)
	}
	return labels
}

func addLabel(tb testing.TB, proj *gogitlab.Project, name, color string) *gogitlab.Label {
	l, _, err := GitLabClient.Labels.CreateLabel(*proj.ID, &gogitlab.CreateLabelOptions{
		Name:  name,
		Color: color,
	})
	if err != nil {
		// The failure happens at wherever we were called, not here
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			tb.Fatalf("Unable to get caller")
		}
		tb.Fatalf("%s:%v %v", path.Base(file), line, err)
	}
	return l
}

func labelsExist(tb testing.TB, proj *gogitlab.Project, expected []*gogitlab.Label) {
	labels := getLabels(tb, *proj.ID)
	for _, exp := range expected {
		found := false
		for _, l := range labels {
			e := *exp
			if exp.Color == "" {
				e.Color = l.Color
			}
			if reflect.DeepEqual(&e, l) {
				found = true
				break
			}
		}
		if !found {
			tb.Fatalf("label %v doesn't exist in %v", exp, labels)
		}
	}
}
