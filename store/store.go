package store

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Load(loc string) ([]byte, error) {
	u, err := url.Parse(loc)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" {
		b, err := ioutil.ReadFile(loc)
		if err != nil {
			return nil, err
		}
		return b, nil
	} else {
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
	//return nil, fmt.Errorf("Unsupported URI to load")
}

func Save(loc string, b []byte, m os.FileMode) error {
	u, err := url.Parse(loc)
	if err != nil {
		return err
	}
	if u.Scheme == "" {
		return ioutil.WriteFile(loc, b, m)
	} else if strings.HasSuffix(u.Host, ".blob.core.windows.net") {
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
	return fmt.Errorf("Unsupported URI to save")
}
