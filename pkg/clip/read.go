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
	"encoding/xml"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
)

var filenameRegex = regexp.MustCompile(`([[:alnum:]]{6}).XML`)

// ReadXML takes a filename and returns the clip Info structre.
func ReadXML(filename string, clipXML []byte) (*Info, error) {
	var infoStruct clipFile
	err := xml.Unmarshal(clipXML, &infoStruct)
	if err != nil {
		return nil, err
	}

	// Get filenames for the clip's component files.
	clipDirname, clipFilename := filepath.Split(filename)
	// We want the CONTENTS folder, so go one directory up.
	clipDirname = filepath.Dir(filepath.Clean(clipDirname))
	groups := filenameRegex.FindStringSubmatch(clipFilename)
	clipPrefix := groups[1]
	for i := range infoStruct.Info.Audio {
		fn := fmt.Sprintf("%s%02d.%s", clipPrefix, i, infoStruct.Info.Audio[i].Format)
		infoStruct.Info.Audio[i].Filename = filepath.Join(clipDirname, "AUDIO", fn)
	}

	// If there's more than one video, this will not work. At all.
	if len(infoStruct.Info.Video) != 1 {
		log.Fatalf("I can't handle multiple video tracks in one clip yet.")
	}

	for i := range infoStruct.Info.Video {
		fn := clipPrefix + "." + infoStruct.Info.Video[i].Format
		infoStruct.Info.Video[i].Filename = filepath.Join(clipDirname, "VIDEO", fn)
	}
	infoStruct.Info.OutputFile = clipPrefix + ".MKV"
	infoStruct.Info.OutputDir = clipDirname

	return &infoStruct.Info, err
}
