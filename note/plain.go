// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package note

// Plain is a "plain old Go object" representation of a note that uses pointers
// to refer directly to related notes in memory.
//
// Most of package notes is designed for working with entire note maps which
// may not fit in memory. Plain is better suited for working with limited
// subgraphs, for example when decoding subgraphs that are to be merged into a
// note map.
type Plain struct {
	ID          ID
	ValueString string
	ValueType   *Plain
	Contents    []*Plain
	Types       []*Plain
}

// GraphNote returns a proxy to x that implements the GraphNote interface.
func (x *Plain) GraphNote() GraphNote {
	return graphPlain{x}
}

type graphPlain struct{ *Plain }

func (x graphPlain) GetID() ID { return x.ID }
func (x graphPlain) GetValue() (string, GraphNote, error) {
	var vt GraphNote
	if x.ValueType == nil {
		vt = EmptyNote(EmptyID)
	} else {
		vt = x.ValueType.GraphNote()
	}
	return x.ValueString, vt, nil
}
func (x graphPlain) GetContents() ([]GraphNote, error) {
	return nmslice(x.Contents)
}
func (x graphPlain) GetTypes() ([]GraphNote, error) {
	return nmslice(x.Types)
}

func nmslice(ps []*Plain) ([]GraphNote, error) {
	gs := make([]GraphNote, len(ps))
	for i, p := range ps {
		gs[i] = p.GraphNote()
	}
	return gs, nil
}
