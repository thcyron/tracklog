package gpx

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

const nsGPX11 = "http://www.topografix.com/GPX/1/1"

var (
	ErrBadRootTag = errors.New("gpx: root element must be <gpx>")
	ErrGPX11Only  = errors.New("gpx: can only parse GPX 1.1 documents")
)

// Document represents a GPX document.
type Document struct {
	Version  string
	Metadata Metadata
	Tracks   []Track
}

// Distance returns the document’s total distance in meters.
func (d Document) Distance() float64 {
	var distance float64
	for _, t := range d.Tracks {
		distance += t.Distance()
	}
	return distance
}

// Duration returns the document’s total duration.
func (d Document) Duration() time.Duration {
	var distance int64
	for _, t := range d.Tracks {
		distance += int64(t.Duration())
	}
	return time.Duration(distance)
}

// Start returns the start time of the first track.
func (d Document) Start() time.Time {
	if len(d.Tracks) == 0 {
		return time.Time{}
	}
	return d.Tracks[0].Start()
}

// End returns the end time of the last track.
func (d Document) End() time.Time {
	if len(d.Tracks) == 0 {
		return time.Time{}
	}
	return d.Tracks[len(d.Tracks)-1].End()
}

// Metadata provides additional information about a GPX document.
type Metadata struct {
	Name        string
	Description string
	Author      Person
	Copyright   Copyright
	Link        Link
	Time        time.Time
	Keywords    string
	Bounds      Bounds
	Extensions  []xml.Token
}

// Person represents a person.
type Person struct {
	Name  string
	Email Email
	Link  Link
}

// Email represents an email address.
type Email struct {
	ID     string
	Domain string
}

// Link represents a link.
type Link struct {
	Href string
	Text string
	Type string
}

// Copyright provides information about the copyright holder and license.
type Copyright struct {
	Author  string
	Year    int
	License string
}

// Bounds provides information about the document bounds.
type Bounds struct {
	MinLatitude  float64
	MinLongitude float64
	MaxLatitude  float64
	MaxLongitude float64
}

// Track represents a track.
type Track struct {
	Name     string
	Segments []Segment
}

// Distance returns the track’s total distance in meters.
func (t Track) Distance() float64 {
	var distance float64
	for _, s := range t.Segments {
		distance += s.Distance()
	}
	return distance
}

// Duration returns the track’s total duration.
func (t Track) Duration() time.Duration {
	var distance int64
	for _, s := range t.Segments {
		distance += int64(s.Duration())
	}
	return time.Duration(distance)
}

// Start returns the start time of the first segment.
func (t Track) Start() time.Time {
	if len(t.Segments) == 0 {
		return time.Time{}
	}
	return t.Segments[0].Start()
}

// End returns the end time of the last segment.
func (t Track) End() time.Time {
	if len(t.Segments) == 0 {
		return time.Time{}
	}
	return t.Segments[len(t.Segments)-1].End()
}

// Segments represents a track segment.
type Segment struct {
	Points []Point
}

// Distance returns the segment’s total distance in meters.
func (s Segment) Distance() float64 {
	var distance float64
	for i, p := range s.Points {
		if i > 0 {
			distance += s.Points[i-1].DistanceTo(p)
		}
	}
	return distance
}

// Duration returns the segment’s total duration.
func (s Segment) Duration() time.Duration {
	ln := len(s.Points)
	if ln < 2 {
		return time.Duration(0)
	}
	return s.Points[ln-1].Time.Sub(s.Points[0].Time)
}

// Start returns the time of the first point.
func (s Segment) Start() time.Time {
	if len(s.Points) == 0 {
		return time.Time{}
	}
	return s.Points[0].Time
}

// End returns the time of the last point.
func (s Segment) End() time.Time {
	if len(s.Points) == 0 {
		return time.Time{}
	}
	return s.Points[len(s.Points)-1].Time
}

// Point represents a track point. Extensions contains the raw XML tokens
// of the point’s extensions if it has any (excluding the <extensions>
// start and end tag).
type Point struct {
	Latitude   float64
	Longitude  float64
	Elevation  float64
	Time       time.Time
	Extensions []xml.Token
}

// DistanceTo returns the distance in meters to point p2.
func (p Point) DistanceTo(p2 Point) float64 {
	return haversine(p.Latitude, p.Longitude, p2.Latitude, p2.Longitude)
}

// Decoder decodes a GPX document from an input stream.
type Decoder struct {
	Strict bool
	r      io.Reader
	ts     tokenStream
}

// NewDecoder creates a new decoder reading from r. The decoder
// operates in strict mode.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		Strict: true,
		r:      r,
	}
}

// Decode decodes a document.
func (d *Decoder) Decode() (doc Document, err error) {
	dec := xml.NewDecoder(d.r)
	d.ts = tokenStream{dec}

	se, err := d.findGPX()
	if err != nil {
		return doc, err
	}

	return d.consumeGPX(se)
}

func (d *Decoder) findGPX() (se xml.StartElement, err error) {
	for {
		tok, err := d.ts.Token()
		if err != nil {
			return se, err
		}
		if se, ok := tok.(xml.StartElement); ok {
			if se.Name.Local != "gpx" {
				return se, ErrBadRootTag
			}
			if se.Name.Space != nsGPX11 {
				return se, ErrGPX11Only
			}
			return se, nil
		}
	}
}

func (d *Decoder) consumeGPX(se xml.StartElement) (doc Document, err error) {
	for _, a := range se.Attr {
		switch a.Name.Local {
		case "version":
			doc.Version = a.Value
		}
	}

	for {
		tok, err := d.ts.Token()
		if err != nil {
			return doc, err
		}
		switch tok.(type) {
		case xml.StartElement:
			se := tok.(xml.StartElement)
			switch se.Name.Local {
			case "trk":
				track, err := d.consumeTrack(se)
				if err != nil {
					return doc, err
				}
				doc.Tracks = append(doc.Tracks, track)
			case "metadata":
				metadata, err := d.consumeMetadata(se)
				if err != nil {
					return doc, err
				}
				doc.Metadata = metadata
			default:
				if err := d.ts.skipTag(); err != nil {
					return doc, err
				}
			}
		case xml.EndElement:
			return doc, nil
		}
	}
}

func (d *Decoder) consumeMetadata(se xml.StartElement) (metadata Metadata, err error) {
	for {
		tok, err := d.ts.Token()
		if err != nil {
			return metadata, err
		}
		switch tok.(type) {
		case xml.StartElement:
			se := tok.(xml.StartElement)
			switch se.Name.Local {
			case "time":
				t, err := d.ts.consumeTime()
				if err != nil {
					return metadata, err
				}
				metadata.Time = t
			case "name":
				s, err := d.ts.consumeString()
				if err != nil {
					return metadata, err
				}
				metadata.Name = s
			case "desc":
				s, err := d.ts.consumeString()
				if err != nil {
					return metadata, err
				}
				metadata.Description = s
			case "link":
				link, err := d.consumeLink(se)
				if err != nil {
					return metadata, err
				}
				metadata.Link = link
			case "keywords":
				s, err := d.ts.consumeString()
				if err != nil {
					return metadata, err
				}
				metadata.Keywords = s
			case "copyright":
				copyright, err := d.consumeCopyright(se)
				if err != nil {
					return metadata, err
				}
				metadata.Copyright = copyright
			case "bounds":
				bounds, err := d.consumeBounds(se)
				if err != nil {
					return metadata, err
				}
				metadata.Bounds = bounds
			case "author":
				person, err := d.consumePerson(se)
				if err != nil {
					return metadata, err
				}
				metadata.Author = person
			default:
				if err := d.ts.skipTag(); err != nil {
					return metadata, err
				}
			}
		case xml.EndElement:
			return metadata, nil
		}
	}
}

func (d *Decoder) consumeLink(se xml.StartElement) (link Link, err error) {
	for _, a := range se.Attr {
		switch a.Name.Local {
		case "href":
			link.Href = a.Value
		}
	}

	for {
		tok, err := d.ts.Token()
		if err != nil {
			return link, err
		}
		switch tok.(type) {
		case xml.StartElement:
			se := tok.(xml.StartElement)
			switch se.Name.Local {
			case "text":
				s, err := d.ts.consumeString()
				if err != nil {
					return link, err
				}
				link.Text = s
			case "type":
				s, err := d.ts.consumeString()
				if err != nil {
					return link, err
				}
				link.Type = s
			default:
				if err := d.ts.skipTag(); err != nil {
					return link, err
				}
			}
		case xml.EndElement:
			return link, nil
		}
	}
}

func (d *Decoder) consumeCopyright(se xml.StartElement) (copyright Copyright, err error) {
	for _, a := range se.Attr {
		switch a.Name.Local {
		case "author":
			copyright.Author = a.Value
		}
	}

	for {
		tok, err := d.ts.Token()
		if err != nil {
			return copyright, err
		}
		switch tok.(type) {
		case xml.StartElement:
			se := tok.(xml.StartElement)
			switch se.Name.Local {
			case "year":
				i, err := d.ts.consumeInt()
				if err != nil {
					return copyright, err
				}
				copyright.Year = i
			case "license":
				s, err := d.ts.consumeString()
				if err != nil {
					return copyright, err
				}
				copyright.License = s
			default:
				if err := d.ts.skipTag(); err != nil {
					return copyright, err
				}
			}
		case xml.EndElement:
			return copyright, nil
		}
	}
}

func (d *Decoder) consumeBounds(se xml.StartElement) (bounds Bounds, err error) {
	for _, a := range se.Attr {
		switch a.Name.Local {
		case "minlat":
			minlat, err := strconv.ParseFloat(a.Value, 64)
			if err == nil {
				bounds.MinLatitude = minlat
			} else if d.Strict {
				return bounds, fmt.Errorf("gpx: invalid <bounds> minlat: %s", err)
			}
		case "maxlat":
			maxlat, err := strconv.ParseFloat(a.Value, 64)
			if err == nil {
				bounds.MaxLatitude = maxlat
			} else if d.Strict {
				return bounds, fmt.Errorf("gpx: invalid <bounds> maxlat: %s", err)
			}
		case "minlon":
			minlon, err := strconv.ParseFloat(a.Value, 64)
			if err == nil {
				bounds.MinLongitude = minlon
			} else if d.Strict {
				return bounds, fmt.Errorf("gpx: invalid <bounds> minlon: %s", err)
			}
		case "maxlon":
			maxlon, err := strconv.ParseFloat(a.Value, 64)
			if err == nil {
				bounds.MaxLongitude = maxlon
			} else if d.Strict {
				return bounds, fmt.Errorf("gpx: invalid <bounds> maxlon: %s", err)
			}
		}
	}

	if d.Strict {
		tok, err := d.ts.Token()
		if err != nil {
			return bounds, err
		}
		if _, ok := tok.(xml.EndElement); !ok {
			return bounds, fmt.Errorf("gpx: <bounds> not a self-closing element")
		}
		return bounds, nil
	}

	for {
		tok, err := d.ts.Token()
		if err != nil {
			return bounds, err
		}
		switch tok.(type) {
		case xml.StartElement:
			if err := d.ts.skipTag(); err != nil {
				return bounds, err
			}
		case xml.EndElement:
			return bounds, nil
		}
	}
}

func (d *Decoder) consumePerson(se xml.StartElement) (person Person, err error) {
	for {
		tok, err := d.ts.Token()
		if err != nil {
			return person, err
		}
		switch tok.(type) {
		case xml.StartElement:
			se := tok.(xml.StartElement)
			switch se.Name.Local {
			case "name":
				name, err := d.ts.consumeString()
				if err != nil {
					return person, err
				}
				person.Name = name
			case "email":
				email, err := d.consumeEmail(se)
				if err != nil {
					return person, err
				}
				person.Email = email
			case "link":
				link, err := d.consumeLink(se)
				if err != nil {
					return person, err
				}
				person.Link = link
			default:
				if err := d.ts.skipTag(); err != nil {
					return person, err
				}
			}
		case xml.EndElement:
			return person, nil
		}
	}
}

func (d *Decoder) consumeEmail(se xml.StartElement) (email Email, err error) {
	for _, a := range se.Attr {
		switch a.Name.Local {
		case "id":
			email.ID = a.Value
		case "domain":
			email.Domain = a.Value
		}
	}

	if d.Strict {
		tok, err := d.ts.Token()
		if err != nil {
			return email, err
		}
		if _, ok := tok.(xml.EndElement); !ok {
			return email, fmt.Errorf("gpx: <email> not a self-closing element")
		}
		return email, nil
	}

	for {
		tok, err := d.ts.Token()
		if err != nil {
			return email, err
		}
		switch tok.(type) {
		case xml.StartElement:
			if err := d.ts.skipTag(); err != nil {
				return email, err
			}
		case xml.EndElement:
			return email, nil
		}
	}
}

func (d *Decoder) consumeTrack(se xml.StartElement) (track Track, err error) {
	for {
		tok, err := d.ts.Token()
		if err != nil {
			return track, err
		}
		switch tok.(type) {
		case xml.StartElement:
			se := tok.(xml.StartElement)
			switch se.Name.Local {
			case "trkseg":
				seg, err := d.consumeSegment(se)
				if err != nil {
					return track, err
				}
				track.Segments = append(track.Segments, seg)
			case "name":
				name, err := d.ts.consumeString()
				if err != nil {
					return track, err
				}
				track.Name = name
			default:
				if err := d.ts.skipTag(); err != nil {
					return track, err
				}
			}
		case xml.EndElement:
			return track, nil
		}
	}
}

func (d *Decoder) consumeSegment(se xml.StartElement) (seg Segment, err error) {
	for {
		tok, err := d.ts.Token()
		if err != nil {
			return seg, err
		}
		switch tok.(type) {
		case xml.StartElement:
			se := tok.(xml.StartElement)
			switch se.Name.Local {
			case "trkpt":
				point, err := d.consumePoint(se)
				if err != nil {
					return seg, err
				}
				seg.Points = append(seg.Points, point)
			default:
				if err := d.ts.skipTag(); err != nil {
					return seg, err
				}
			}
		case xml.EndElement:
			return seg, nil
		}
	}
}

func (d *Decoder) consumePoint(se xml.StartElement) (point Point, err error) {
	for _, a := range se.Attr {
		switch a.Name.Local {
		case "lat":
			lat, err := strconv.ParseFloat(a.Value, 64)
			if err == nil {
				point.Latitude = lat
			} else if d.Strict {
				return point, fmt.Errorf("gpx: invalid <trkpt> lat: %s", err)
			}
		case "lon":
			lon, err := strconv.ParseFloat(a.Value, 64)
			if err == nil {
				point.Longitude = lon
			} else if d.Strict {
				return point, fmt.Errorf("gpx: invalid <trkpt> lon: %s", err)
			}
		}
	}

	for {
		tok, err := d.ts.Token()
		if err != nil {
			return point, err
		}
		switch tok.(type) {
		case xml.StartElement:
			se := tok.(xml.StartElement)
			switch se.Name.Local {
			case "ele":
				ele, err := d.ts.consumeFloat()
				if err != nil {
					return point, err
				}
				point.Elevation = ele
			case "time":
				t, err := d.ts.consumeTime()
				if err != nil {
					return point, err
				}
				point.Time = t
			case "extensions":
				exts, err := d.consumeExtensions(se)
				if err != nil {
					return point, err
				}
				point.Extensions = exts
			default:
				if err := d.ts.skipTag(); err != nil {
					return point, err
				}
			}
		case xml.EndElement:
			return point, nil
		}
	}
}

func (d *Decoder) consumeExtensions(se xml.StartElement) (tokens []xml.Token, err error) {
	lvl := 0

	for {
		tok, err := d.ts.Token()
		if err != nil {
			return tokens, err
		}
		switch tok.(type) {
		case xml.StartElement:
			lvl++
		case xml.EndElement:
			if lvl == 0 {
				return tokens, nil
			}
			lvl--
		}
		tokens = append(tokens, xml.CopyToken(tok))
	}
}
