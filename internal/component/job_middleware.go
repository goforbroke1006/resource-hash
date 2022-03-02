package component

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"

	"resource-hash/domain"
)

func NewCheckLinkJob(link string, out chan<- domain.OutputChunk) func() error {
	return func() error {
		resp, err := http.Get(link)
		if err != nil {
			return err
		}
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		sum := md5.Sum(content)

		out <- domain.OutputChunk{
			Url:  link,
			Hash: fmt.Sprintf("%x", sum),
		}

		return nil
	}
}
