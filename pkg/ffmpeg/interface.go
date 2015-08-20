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

// Graph is an object that produces an ffmpeg filter graph.
// Each function taskes as input either a filename or a node label,
// which will be mapped to the function's input, and returns the node label
// of the function's output.
type Graph interface {
	AddVideo(inputNode string) string
	AddAudio(inputNode string) string
	MergeAudio(inputNode ...string) string
	MixAudio(inputNode ...string) string

	Filter() string

	// Inputs() []string
	// Outputs() []string
}
