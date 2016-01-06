package gpx

import (
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"time"
)

// A tokener provides tokens.
type tokener interface {
	Token() (xml.Token, error)
}

// A sliceTokener provides tokens from a slice.
type sliceTokener struct {
	tokens []xml.Token
}

func (t *sliceTokener) Token() (xml.Token, error) {
	if len(t.tokens) == 0 {
		return nil, io.EOF
	}
	tok := t.tokens[0]
	t.tokens = t.tokens[1:]
	return tok, nil
}

// A tokenStream operates on a stream of tokens from a tokener.
type tokenStream struct {
	tokener
}

func (ts *tokenStream) consumeString() (string, error) {
	var s string
	for {
		tok, err := ts.Token()
		if err != nil {
			return "", err
		}
		switch tok.(type) {
		case xml.CharData:
			s += string(tok.(xml.CharData))
		case xml.EndElement:
			return s, nil
		default:
			return "", errors.New("gpx: unexpected element while reading string")
		}
	}
}

func (ts *tokenStream) consumeTime() (time.Time, error) {
	s, err := ts.consumeString()
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339Nano, s)
}

func (ts *tokenStream) consumeFloat() (float64, error) {
	s, err := ts.consumeString()
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(s, 64)
}

func (ts *tokenStream) consumeInt() (int, error) {
	s, err := ts.consumeString()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

func (ts *tokenStream) skipTag() error {
	for {
		tok, err := ts.Token()
		if err != nil {
			return err
		}
		switch tok.(type) {
		case xml.StartElement:
			if err := ts.skipTag(); err != nil {
				return nil
			}
		case xml.EndElement:
			return nil
		}
	}
}
