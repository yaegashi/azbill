package store

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

type Store struct {
	Dir string
}

func NewStore(dir string) (*Store, error) {
	dir, err := homedir.Expand(dir)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	return &Store{Dir: dir}, nil
}

func (s *Store) Path(path string) string {
	if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "./") {
		return path
	}
	return filepath.Join(s.Dir, path)
}

func (s *Store) Location(loc string, redact bool) string {
	u, err := url.Parse(loc)
	if err != nil {
		return loc
	}
	if redact {
		if u.RawQuery != "" {
			u.RawQuery = "..."
		}
	}
	switch u.Scheme {
	case "file", "":
		if u.Path != "" {
			return s.Path(u.Path)
		}
	}
	return u.String()
}

func (s *Store) ReadFile(loc string) ([]byte, error) {
	loc = s.Location(loc, false)
	u, err := url.Parse(loc)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "":
		b, err := ioutil.ReadFile(loc)
		if err != nil {
			return nil, err
		}
		return b, nil
	case "https":
		res, err := http.Get(loc)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s", res.Status)
		}
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	return nil, fmt.Errorf("Unsupported URI to load")
}

func (s *Store) WriteFile(loc string, b []byte, m os.FileMode) error {
	loc = s.Location(loc, false)
	u, err := url.Parse(loc)
	if err != nil {
		return err
	}
	switch u.Scheme {
	case "":
		if strings.HasPrefix(loc, s.Dir) {
			err = os.MkdirAll(filepath.Dir(s.Dir), 0755)
			if err != nil {
				return err
			}
		}
		return ioutil.WriteFile(loc, b, m)
	case "https":
		if strings.HasSuffix(u.Host, ".blob.core.windows.net") {
			cli := &http.Client{}
			req, err := http.NewRequest(http.MethodPut, loc, bytes.NewBuffer(b))
			if err != nil {
				return err
			}
			req.Header.Set("x-ms-blob-type", "BlockBlob")
			res, err := cli.Do(req)
			if err != nil {
				return err
			}
			defer res.Body.Close()
			if res.StatusCode != http.StatusCreated {
				return fmt.Errorf("%s", res.Status)
			}
			return nil
		}
	}
	return fmt.Errorf("Unsupported URI to save")
}
