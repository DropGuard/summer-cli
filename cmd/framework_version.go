package cmd

// FrameworkVersion is the Summer Framework version new projects default to.
// Injected at release build time alongside the CLI version so a CLI release never
// pins users to a stale framework line by accident; overridable per-create via
// --framework-version.
var FrameworkVersion = "0.3.2"
