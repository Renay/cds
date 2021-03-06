package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ovh/cds/sdk"
)

func runParseJunitTestResultAction(a *sdk.Action, ab sdk.ActionBuild) sdk.Result {
	var res sdk.Result
	res.Status = sdk.StatusFail

	// Retrieve build info
	var proj, app, pip, bnS, envName string
	for _, p := range ab.Args {
		switch p.Name {
		case "cds.pipeline":
			pip = p.Value
			break
		case "cds.project":
			proj = p.Value
			break
		case "cds.application":
			app = p.Value
			break
		case "cds.buildNumber":
			bnS = p.Value
			break
		case "cds.environment":
			envName = p.Value
			break
		}
	}

	var p string
	for _, a := range a.Parameters {
		if a.Name == "path" {
			p = a.Value
			break
		}
	}

	if p == "" {
		sendLog(ab.ID, sdk.JUnitAction, fmt.Sprintf("UnitTest parser: path not provided"))
		return res
	}

	files, err := filepath.Glob(p)
	if err != nil {
		sendLog(ab.ID, sdk.JUnitAction, fmt.Sprintf("UnitTest parser: Cannot find requested files, invalid pattern"))
		return res
	}

	var v sdk.Tests
	for _, f := range files {
		var ftests sdk.Tests

		data, err := ioutil.ReadFile(f)
		if err != nil {
			sendLog(ab.ID, sdk.JUnitAction, fmt.Sprintf("UnitTest parser: cannot read file %s (%s)", f, err))
			return res
		}

		err = xml.Unmarshal([]byte(data), &v)
		if err != nil {
			sendLog(ab.ID, sdk.JUnitAction, fmt.Sprintf("UnitTest parser: cannot interpret file %s (%s)", f, err))
			return res
		}

		// Is it nosetests format ?
		if s, ok := parseNoseTests(data); ok {
			ftests.TestSuites = append(ftests.TestSuites, s)
		}

		v.TestSuites = append(v.TestSuites, ftests.TestSuites...)
	}
	// update global stats
	for _, s := range v.TestSuites {
		v.Total += s.Total
		v.TotalOK += (s.Total - s.Failures)
		v.TotalKO += s.Failures
		v.TotalSkipped += s.Skip
	}

	res.Status = sdk.StatusSuccess
	for _, s := range v.TestSuites {
		if s.Failures > 0 {
			sendLog(ab.ID, sdk.JUnitAction, fmt.Sprintf("JUnit parser: %s has %d failed tests", s.Name, s.Failures))
			res.Status = sdk.StatusFail
		}
	}

	if v.Total == 0 {
		sendLog(ab.ID, sdk.JUnitAction, "JUnit parser: No tests")
		res.Status = sdk.StatusFail
	}

	data, err := json.Marshal(v)
	if err != nil {
		res.Status = sdk.StatusFail
		sendLog(ab.ID, sdk.JUnitAction, fmt.Sprintf("JUnit parse: failed to send tests details: %s", err))
		return res
	}

	uri := fmt.Sprintf("/project/%s/application/%s/pipeline/%s/build/%s/test?envName=%s", proj, app, pip, bnS, envName)
	_, code, err := sdk.Request("POST", uri, data)
	if err == nil && code > 300 {
		err = fmt.Errorf("HTTP %d", code)
	}
	if err != nil {
		res.Status = sdk.StatusFail
		sendLog(ab.ID, sdk.JUnitAction, fmt.Sprintf("JUnit parse: failed to send tests details: %s", err))
		return res
	}

	return res
}

func parseNoseTests(data []byte) (sdk.TestSuite, bool) {
	var s sdk.TestSuite
	err := xml.Unmarshal([]byte(data), &s)
	if err != nil {
		return s, false
	}

	if s.Name == "" {
		return s, false
	}

	return s, true
}
