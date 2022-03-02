package component

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"

	"resource-hash/domain"
)

// NewCheckLinkJob use closure to prepare job that would send reports with hashes
func NewCheckLinkJob(link string, out chan<- domain.OutputChunk) func() error {
	return func() error {
		chunk := domain.OutputChunk{
			Url: link,
		}
		defer func() {
			out <- chunk
		}()

		resp, err := http.Get(link)
		if err != nil {
			chunk.Warn = err
			return nil
		}
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			chunk.Warn = err
			return nil
		}

		sum := md5.Sum(content)
		chunk.Hash = fmt.Sprintf("%x", sum)

		return nil
	}
}
