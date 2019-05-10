package junitxml

// Document represents a JUnitXML document based on the schema.
type Document struct {
	// Suites represents the top-level testsuites element of the file.
	Suites []JUnitTestSuites
}
