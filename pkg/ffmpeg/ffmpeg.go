/*
Copyright © 2015 Brad Ackerman.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package ffmpeg

import (
	"fmt"
	"strings"
)

// FIXME: Need to implement graph consistency.

type filterGraph struct {
	// Every node.
	nodes []string
	// The inputs.
	inputs []string
	// The outputs.
	outputs []string
	// Counter for number of video nodes
	vCounter int
	// Counter for number of audio nodes
	aCounter int
}

// NewFilterGraph creates a filter graph, joy.
func NewFilterGraph() Graph {
	return &filterGraph{}
}

func (g *filterGraph) AddVideo(inputNode string) string {
	g.vCounter++
	label := fmt.Sprintf("[v%d]", g.vCounter)
	node := fmt.Sprintf("movie='%s' %s", inputNode, label)
	g.nodes = append(g.nodes, node)
	return label
}

func (g *filterGraph) AddAudio(inputNode string) string {
	g.aCounter++
	label := fmt.Sprintf("[a%d]", g.aCounter)
	node := fmt.Sprintf("amovie='%s' %s", inputNode, label)
	g.nodes = append(g.nodes, node)
	return label
}

func (g *filterGraph) MergeAudio(inputNodes ...string) string {
	inputs := strings.Join(inputNodes, " ")
	g.aCounter++
	outputLabel := fmt.Sprintf("[a%d]", g.aCounter)
	node := fmt.Sprintf("%s amerge=inputs=%d %s", inputs, len(inputNodes), outputLabel)
	g.nodes = append(g.nodes, node)
	return outputLabel
}

func (g *filterGraph) MixAudio(inputNodes ...string) string {
	inputs := strings.Join(inputNodes, " ")
	g.aCounter++
	outputLabel := fmt.Sprintf("[a%d]", g.aCounter)
	node := fmt.Sprintf("%s amix=inputs=%d %s", inputs, len(inputNodes), outputLabel)
	g.nodes = append(g.nodes, node)
	return outputLabel
}

// func (g *filterGraph) Inputs() []string {
// 	return []string{}
// }
// func (g *filterGraph) Outputs() []string {
// 	return []string{}
// }

func (g *filterGraph) Filter() string {
	return strings.Join(g.nodes, ";\n")
}
