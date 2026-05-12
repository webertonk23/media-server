package utils
import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
)
type FFProbeOutput struct {
	Streams  []FFStream  `json:"streams"`
	Format   FFFormat    `json:"format"`
	Chapters []FFChapter `json:"chapters"`
}
type FFChapter struct {
	ID        int     `json:"id"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	Tags      Tags    `json:"tags,omitempty"`
}
type FFStream struct {
	Index         int    `json:"index"`
	CodecName     string `json:"codec_name"`
	CodecType     string `json:"codec_type"`
	Width         int    `json:"width,omitempty"`
	Height        int    `json:"height,omitempty"`
	BitRate       string `json:"bit_rate"`
	Tags          Tags   `json:"tags,omitempty"`
	ChannelLayout string `json:"channel_layout,omitempty"`
}
type FFFormat struct {
	Filename string `json:"filename"`
	Duration string `json:"duration"`
	Size     string `json:"size"`
	Tags     Tags   `json:"tags,omitempty"`
}
type Tags struct {
	Title    string `json:"title,omitempty"`
	Language string `json:"language,omitempty"`
	Season   string `json:"season_number,omitempty"`
	Episode  string `json:"episode_number,omitempty"`
	Show     string `json:"show,omitempty"`
}
func Probe(path string) (*FFProbeOutput, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		"-show_chapters",
		path,
	)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe error: %v", err)
	}
	var result FFProbeOutput
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %v", err)
	}
	return &result, nil
}
func (f *FFProbeOutput) GetVideoPath() string {
	return f.Format.Filename
}
func (f *FFProbeOutput) HasH264() bool {
	for _, s := range f.Streams {
		if s.CodecType == "video" && s.CodecName == "h264" {
			return true
		}
	}
	return false
}
func (f *FFProbeOutput) HasCompatibleAudio() bool {
	for _, s := range f.Streams {
		if s.CodecType == "audio" && (s.CodecName == "aac" || s.CodecName == "mp3" || s.CodecName == "ac3") {
			return true
		}
	}
	return false
}
func (f *FFProbeOutput) GetAudioStreams() []FFStream {
	var audios []FFStream
	for _, s := range f.Streams {
		if s.CodecType == "audio" {
			audios = append(audios, s)
		}
	}
	return audios
}
func (f *FFProbeOutput) GetSubtitles() []FFStream {
	var subs []FFStream
	for _, s := range f.Streams {
		if s.CodecType == "subtitle" {
			subs = append(subs, s)
		}
	}
	return subs
}
func (f *FFProbeOutput) GetMetadata() Tags {
	tags := f.Format.Tags
	for _, s := range f.Streams {
		if tags.Title == "" && s.Tags.Title != "" {
			tags.Title = s.Tags.Title
		}
		if tags.Show == "" && s.Tags.Show != "" {
			tags.Show = s.Tags.Show
		}
	}
	return tags
}
func (c *FFChapter) GetStartTime() float64 {
	val, _ := strconv.ParseFloat(c.StartTime, 64)
	return val
}
func (c *FFChapter) GetEndTime() float64 {
	val, _ := strconv.ParseFloat(c.EndTime, 64)
	return val
}
