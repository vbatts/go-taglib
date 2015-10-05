package taglib

/*
#include <stdlib.h>
#include <taglib/tag_c.h>
#cgo LDFLAGS: -ltag_c
*/
import "C"
import "unsafe"

// Tags are the tag fields accesible
type Tags struct {
	Title, Artist, Album, Comment, Genre string
	Year, Track                          int
}

// Properties are the attributes of the File
type Properties struct {
	Length, Bitrate, Samplerate, Channels int
}

// File is the means to access properties, and set tags.
type File C.TagLib_File

// Open returns a taglib.File from music file filename. If nil, then the file
// could not be opened.
func Open(filename string) *File {
	fp := C.CString(filename)
	defer C.free(unsafe.Pointer(fp))
	return (*File)(C.taglib_file_new(fp))
}

var cTrue = C.int(1)

// Valid returns true if the file is open and readble and valid information for
// the Tag and / or AudioProperties was found.
func (f *File) Valid() bool {
	return C.taglib_file_is_valid((*C.TagLib_File)(f)) == cTrue
}

// Save the file to disk.
func (f *File) Save() bool {
	return C.taglib_file_save((*C.TagLib_File)(f)) == cTrue
}

// Close and free a taglib.File
func (f *File) Close() {
	C.taglib_file_free((*C.TagLib_File)(f))
}

// GetTags returns the ID3 taglib.Tags from this taglib.File
func (f *File) GetTags() *Tags {
	ts := C.taglib_file_tag((*C.TagLib_File)(f))

	a := Tags{}

	if ts != nil {
		a.Title = C.GoString(C.taglib_tag_title(ts))
		a.Artist = C.GoString(C.taglib_tag_artist(ts))
		a.Album = C.GoString(C.taglib_tag_album(ts))
		a.Comment = C.GoString(C.taglib_tag_comment(ts))
		a.Genre = C.GoString(C.taglib_tag_genre(ts))
		a.Year = int(C.taglib_tag_year(ts))
		a.Track = int(C.taglib_tag_track(ts))
	}

	defer C.taglib_tag_free_strings()
	return &a
}

// GetProperties returns the taglib.Properties from this taglib.File
func (f *File) GetProperties() *Properties {
	ap := C.taglib_file_audioproperties((*C.TagLib_File)(f))
	if ap == nil {
		return nil
	}

	p := Properties{}
	p.Length = int(C.taglib_audioproperties_length(ap))
	p.Bitrate = int(C.taglib_audioproperties_bitrate(ap))
	p.Samplerate = int(C.taglib_audioproperties_samplerate(ap))
	p.Channels = int(C.taglib_audioproperties_channels(ap))

	defer C.taglib_tag_free_strings()
	return &p
}

// GetTags returns the ID3 taglib.Tags from filename
func GetTags(filename string) *Tags {
	tf := Open(filename)
	if tf == nil {
		return nil
	}
	defer tf.Close()
	return tf.GetTags()
}

// GetProperties returs the taglib.Properties from filename
func GetProperties(filename string) *Properties {
	tf := Open(filename)
	if tf == nil {
		return nil
	}
	defer tf.Close()
	return tf.GetProperties()
}
