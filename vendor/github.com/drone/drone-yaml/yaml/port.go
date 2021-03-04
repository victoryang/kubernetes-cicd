// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package yaml

type (

	// Port represents a network port in a single container.
	Port struct {
		Port     int    `json:"port,omitempty"`
		Host     int    `json:"host,omitempty"`
		Protocol string `json:"protocol,omitempty"`
	}

	port struct {
		Port     int
		Host     int
		Protocol string
	}
)

// UnmarshalYAML implements yaml unmarshalling.
func (p *Port) UnmarshalYAML(unmarshal func(interface{}) error) error {
	out := new(port)
	err := unmarshal(&out.Port)
	if err != nil {
		err = unmarshal(&out)
	}
	p.Port = out.Port
	p.Host = out.Host
	p.Protocol = out.Protocol
	return err
}
