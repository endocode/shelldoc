package junitxml

// SPDX-License-Identifier:	MIT
// Copyright (c) 2012 Joel Stemmer
//
// The types in this files are based on an excerpt from https://github.com/jstemmer/go-junit-report.

import (
	"encoding/xml"
	"fmt"
	"time"
)

// JUnitTestSuites is a collection of JUnit test suites.
type JUnitTestSuites struct {
	XMLName xml.Name `xml:"testsuites"`
	Suites  []JUnitTestSuite
}

// JUnitTestSuite is a single JUnit test suite which may contain many
// testcases.
type JUnitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Time       string          `xml:"time,attr"`
	Name       string          `xml:"name,attr"`
	Properties []JUnitProperty `xml:"properties>property,omitempty"`
	TestCases  []JUnitTestCase
}

// JUnitTestCase is a single test case with its result.
type JUnitTestCase struct {
	XMLName     xml.Name          `xml:"testcase"`
	Classname   string            `xml:"classname,attr"`
	Name        string            `xml:"name,attr"`
	Time        string            `xml:"time,attr"`
	SkipMessage *JUnitSkipMessage `xml:"skipped,omitempty"`
	Failure     *JUnitFailure     `xml:"failure,omitempty"`
}

// JUnitSkipMessage contains the reason why a testcase was skipped.
type JUnitSkipMessage struct {
	Message string `xml:"message,attr"`
}

// JUnitProperty represents a key/value pair used to define properties.
type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// JUnitFailure contains data related to a failed test.
type JUnitFailure struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",chardata"`
}

// FormatTime creates a representation of time.Duration as expected in the JUnixXML output.
func FormatTime(d time.Duration) string {
	return fmt.Sprintf("%.3f", d.Seconds())
}

// FormatBenchmarkTime creates a representation of time.Duration as expected in the JUnixXML output for benchmarks.
func FormatBenchmarkTime(d time.Duration) string {
	return fmt.Sprintf("%.9f", d.Seconds())
}

// AddProperty adds a property to the properties of the test suite.
func (suite *JUnitTestSuite) AddProperty(key, value string) {
	prop := JUnitProperty{key, value}
	suite.Properties = append(suite.Properties, prop)
}

// TestCount returns the number of test cases in the test suite.
func (suite *JUnitTestSuite) TestCount() int {
	return len(suite.TestCases)
}

// SuccessCount returns the number of successfully executed test cases in the test suite.
func (suite *JUnitTestSuite) SuccessCount() int {
	counter := 0
	for _, testcase := range suite.TestCases {
		if testcase.Failure == nil {
			counter++
		}
	}
	return counter
}

// FailureCount returns the number of executed test cases in the test suite that have failed.
func (suite *JUnitTestSuite) FailureCount() int {
	return suite.TestCount() - suite.SuccessCount()
}

// RegisterTestCase registers a test case with the test suite. The test count increments.
func (suite *JUnitTestSuite) RegisterTestCase(testcase JUnitTestCase) {
	suite.Tests++
	suite.TestCases = append(suite.TestCases, testcase)
	if suite.Tests != suite.TestCount() {
		panic(fmt.Sprintf("internal constraint violated - Tests and TestCases mismatch: %v", suite))
	}
}
