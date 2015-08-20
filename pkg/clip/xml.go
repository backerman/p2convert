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
	"time"
)

type clipFile struct {
	Info Info `xml:"ClipContent"`
}

// Info is the metadata for a clip as desribed by the CLIP/*.XML files.
type Info struct {
	Name     string         `xml:"ClipName"`
	GlobalID string         `xml:"GlobalClipID"`
	Duration int            `xml:"Duration"`
	EditUnit string         `xml:"EditUnit"` // I don't actually know what this is.
	Video    []VideoEssence `xml:"EssenceList>Video"`
	Audio    []AudioEssence `xml:"EssenceList>Audio"`
	Metadata Metadata       `xml:"ClipMetadata"`
	// Added by us
	OutputFile string
	OutputDir  string // default; can be changed
}

// Metadata provides the clip's metadata.
type Metadata struct {
	Name         string `xml:"UserClipName"`
	DataSource   string
	CreationDate Time       `xml:"Access>CreationDate"`
	LastUpdate   Time       `xml:"Access>LastUpdateDate"`
	DeviceInfo   DeviceInfo `xml:"Device"`
	Start        Time       `xml:"Shoot>StartDate"`
	End          Time       `xml:"Shoot>EndDate"`
	Thumbnail    ThumbnailInfo
}

// VideoEssence provides metadata about the video in a clip.
type VideoEssence struct {
	Format    string `xml:"VideoFormat"`
	Codec     string
	FrameRate string
	// DropFrame        bool `xml:"FrameRate>DropFrameFlag,attr"`
	StartTimecode    string
	StartBinaryGroup string
	AspectRatio      string
	Index            EssenceIndex `xml:"VideoIndex"`

	// filename inserted by us
	Filename string
}

// AudioEssence provides metadata about the audio in a clip.
type AudioEssence struct {
	Format        string `xml:"AudioFormat"`
	SamplingRate  int
	BitsPerSample int
	Index         EssenceIndex `xml:"AudioIndex"`

	// filename inserted by us
	Filename string
}

// EssenceIndex stores the start position and length of an essence.
type EssenceIndex struct {
	StartByteOffset int
	DataSize        int
}

// DeviceInfo stores the type and serial number of the device that produced this
// clip.
type DeviceInfo struct {
	Manufacturer string
	Serial       string `xml:"SerialNo."`
	Model        string `xml:"ModelName"`
}

// ThumbnailInfo ...
type ThumbnailInfo struct {
	FrameOffset int
	Format      string `xml:"ThumbnailFormat"`
	Width       int
	Height      int
}

// Time is a wrapper for time.Time that can be unmarshalled from XML.
type Time struct {
	time.Time
}

const timestampLayout = "2006-01-02T15:04:05-07:00"

// UnmarshalXML does what it says on the tin: unmarshals XML into a time.Time.
func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var timeStr string
	err := d.DecodeElement(&timeStr, &start)
	if err != nil {
		return err
	}
	goTime, err := time.Parse(timestampLayout, timeStr)
	if err != nil {
		return err
	}
	*t = Time{goTime}
	return nil
}
