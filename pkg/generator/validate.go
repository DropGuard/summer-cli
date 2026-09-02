package generator

import (
	"fmt"
	"strings"
	"unicode"
)

// javaKeywords that are illegal as a (package/class) identifier segment.
var javaKeywords = map[string]bool{
	"abstract": true, "assert": true, "boolean": true, "break": true, "byte": true,
	"case": true, "catch": true, "char": true, "class": true, "const": true,
	"continue": true, "default": true, "do": true, "double": true, "else": true,
	"enum": true, "extends": true, "final": true, "finally": true, "float": true,
	"for": true, "goto": true, "if": true, "implements": true, "import": true,
	"instanceof": true, "int": true, "interface": true, "long": true, "native": true,
	"new": true, "package": true, "private": true, "protected": true, "public": true,
	"return": true, "short": true, "static": true, "strictfp": true, "super": true,
	"switch": true, "synchronized": true, "this": true, "throw": true, "throws": true,
	"transient": true, "try": true, "void": true, "volatile": true, "while": true,
	"var": true, "record": true, "yield": true, "_": true,
}

// IsJavaIdentifier reports whether s is a legal Java identifier segment
// (letter/underscore start, then letters/digits/underscores, not a keyword).
func IsJavaIdentifier(s string) bool {
	if s == "" || javaKeywords[s] {
		return false
	}
	for i, r := range s {
		if !unicode.IsLetter(r) && r != '_' && r != '$' && (i == 0 || !unicode.IsDigit(r)) {
			return false
		}
	}
	return true
}

// SanitizePackageSegment derives a legal Java package segment from a free-form
// project/artifact name: illegal characters are dropped ("my-app" -> "myapp").
func SanitizePackageSegment(s string) string {
	var b strings.Builder
	for i, r := range strings.TrimSpace(strings.ToLower(s)) {
		switch {
		case unicode.IsLetter(r) || r == '_':
			b.WriteRune(r)
		case unicode.IsDigit(r) && i > 0:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ValidateGroupId checks every dot-separated segment is a legal Java identifier.
func ValidateGroupId(groupId string) error {
	if groupId == "" {
		return fmt.Errorf("groupId must not be empty")
	}
	for _, seg := range strings.Split(groupId, ".") {
		if !IsJavaIdentifier(seg) {
			return fmt.Errorf("invalid groupId %q: segment %q is not a legal Java identifier", groupId, seg)
		}
	}
	return nil
}

// ValidatePackage validates a full base package (dot-separated identifiers).
func ValidatePackage(pkg string) error {
	if pkg == "" {
		return fmt.Errorf("base package must not be empty")
	}
	for _, seg := range strings.Split(pkg, ".") {
		if !IsJavaIdentifier(seg) {
			return fmt.Errorf("invalid base package %q: segment %q is not a legal Java identifier (hint: symbols are stripped only when DERIVED; pass -p explicitly)", pkg, seg)
		}
	}
	return nil
}

// ValidateProjectName checks the scaffold target name. It becomes a directory
// in the current working directory AND part of Maven coordinates, so path
// separators (traversal like "../evil"), whitespace and shell-hostile
// characters are rejected up front. A name that sanitizes to no legal package
// segment (e.g. "123") is rejected with a pointer to the explicit -p flag.
func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name must not be empty")
	}
	if name == "." || name == ".." {
		return fmt.Errorf("project name %q is not a usable directory name", name)
	}
	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("project name %q must not contain path separators — it is created as a directory inside the current one", name)
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' && r != '.' {
			return fmt.Errorf("project name %q contains unsupported character %q (allowed: letters, digits, '-', '_', '.')", name, string(r))
		}
	}
	if !IsJavaIdentifier(SanitizePackageSegment(name)) {
		return fmt.Errorf("project name %q cannot derive a package segment — pass an explicit one with -p", name)
	}
	return nil
}

// ValidateArtifactId checks the Maven artifactId: no whitespace or path
// separators (it is interpolated into pom.xml and used as a directory name by
// Maven itself).
func ValidateArtifactId(artifactId string) error {
	if artifactId == "" {
		return fmt.Errorf("artifactId must not be empty")
	}
	if strings.ContainsAny(artifactId, " /\\\t\n") {
		return fmt.Errorf("invalid artifactId %q: whitespace and path separators are not allowed", artifactId)
	}
	return nil
}
