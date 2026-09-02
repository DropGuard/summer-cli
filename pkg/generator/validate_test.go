package generator

import "testing"

func TestValidateProjectNameRejectsTraversal(t *testing.T) {
	for _, bad := range []string{"../evil", "a/b", "a\\b", "", ".", "..", "my app", "123"} {
		if err := ValidateProjectName(bad); err == nil {
			t.Fatalf("ValidateProjectName(%q) = nil, want error", bad)
		}
	}
	for _, good := range []string{"my-app", "my_app", "app01", "App.Name"} {
		if err := ValidateProjectName(good); err != nil {
			t.Fatalf("ValidateProjectName(%q) = %v, want nil", good, err)
		}
	}
}

func TestValidateArtifactIdRejectsWhitespaceAndSeparators(t *testing.T) {
	for _, bad := range []string{"", "a b", "a/b", "a\\b", "a\tb"} {
		if err := ValidateArtifactId(bad); err == nil {
			t.Fatalf("ValidateArtifactId(%q) = nil, want error", bad)
		}
	}
	if err := ValidateArtifactId("my-app"); err != nil {
		t.Fatalf("ValidateArtifactId(my-app) = %v, want nil", err)
	}
}
