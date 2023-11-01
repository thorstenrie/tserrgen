// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tserrgen

type tserrconfig struct {
	Root errorsconfig `json:"tserr"` // Root element
}

type errorsconfig struct {
	Path   string   `json:"path"`    // Path to tserr package
	Ver    string   `json:"version"` // Version of tserr package
	Errors []errmsg `json:"errors"`  // Errors
}

type errmsg struct {
	Name    string  `json:"name"`    // Error name
	Comment string  `json:"comment"` // Comment
	Code    string  `json:"code"`    // Error code (HTTP status code from Go standard package http)
	Msg     string  `json:"message"` // Error messages (may contain verbs)
	Param   []param `json:"param"`   // Error parameters
}

type param struct {
	Name    string `json:"name"`    // Parameter name
	Comment string `json:"comment"` // Parameter comment
	Type    string `json:"type"`    // Parameter type
}
