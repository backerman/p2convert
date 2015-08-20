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

package clip

import (
	"path/filepath"

	"github.com/backerman/p2convert/pkg/ffmpeg"
)

// ComplexFilter is an ffmpeg complex filter and its designated output labels.
type ComplexFilter struct {
	Filter      string
	OutputVideo string
	OutputAudio string
}

// GetFilter makes an ffmpeg filter from the clip.
func (clip *Info) GetFilter() ComplexFilter {
	graph := ffmpeg.NewFilterGraph()

	audioTracks := []string{}
	for _, afile := range clip.Audio {
		audioTracks = append(audioTracks, graph.AddAudio(afile.Filename))
	}
	outputAudio := graph.MixAudio(audioTracks...)
	outputVideo := graph.AddVideo(clip.Video[0].Filename)
	filter := graph.Filter()

	return ComplexFilter{
		Filter:      filter,
		OutputVideo: outputVideo,
		OutputAudio: outputAudio,
	}
}

// CommandLine generates an ffmpeg command line to execute the specified filter.
func (clip *Info) CommandLine(filter ComplexFilter) []string {
	// Put the filter into a temporary file.

	cmd := []string{
		// The filter graph
		"-filter_complex", filter.Filter,
		// Map output video
		"-map", filter.OutputVideo,
		"-vcodec", "libx264",
		"-preset", "veryfast",
		// Map output audio
		"-map", filter.OutputAudio,
		"-acodec", "aac", "-strict", "-2",
		"-q:audio", "100",
		"-f", "matroska",
		filepath.Join(clip.OutputDir, clip.OutputFile),
	}

	return cmd
}
