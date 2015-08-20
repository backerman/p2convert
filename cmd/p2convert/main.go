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

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/backerman/p2convert/pkg/clip"
	"github.com/spf13/cobra"
)

/*
Statement of work:

1. Take CONTENTS directory as input.
2. Accept options from user:
- 4x mono? 2x stereo? 1x mono or stereo?
- scale down?
- choice of H264 profiles (fast/slow/etc)
3. Parse the XML files in CLIP directory. Use Adobe code to map codec to input width and height.
4. Using those files, generate the appropriate FFmpeg command to transcode
input to H.264 / AAC in a Matroška container suitable for Youtube upload.
5. Call ffmpeg.
6. Profit!
*/

func runMe(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: Please specify at least one P2 .XML clip file.")
		os.Exit(1)
	}
	for _, filename := range args {
		xmlBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalf("Error! %v", err)
		}
		clip, err := clip.ReadXML(filename, xmlBytes)
		if err != nil {
			log.Fatalf("Error! %v", err)
		}
		if cmd.Flag("outputdir").Value.String() != "" {
			clip.OutputDir = cmd.Flag("outputdir").Value.String()
		}
		filter := clip.GetFilter()
		cmdline := clip.CommandLine(filter)
		for i := range cmdline {
			// Double up the backslashes. This shouldn't be necessary, but FFmpeg
			// seems to require it on Windows.
			cmdline[i] = strings.Replace(cmdline[i], "\\", "\\\\", -1)
		}

		ffCmd := exec.Command("ffmpeg", cmdline...)

		ffCmd.Stdin = os.Stdin
		ffCmd.Stdout = os.Stdout
		ffCmd.Stderr = os.Stderr
		err = ffCmd.Run()
		if err != nil {
			log.Fatalf("Got error running ffmpeg: %v", err)
		}

	}
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "p2convert",
		Short: "Convert a P2 video to a YouTube-uploadble mkv",
		Long: `
		p2convert takes as input a P2 video .XML file, and generates a .mkv file
		that can be uploaded to YouTube or viewed locally.
		`,
		Run: runMe,
	}
	rootCmd.Flags().StringP("outputdir", "d", "", "Directory in which to place output files.")
	rootCmd.Execute()
}
