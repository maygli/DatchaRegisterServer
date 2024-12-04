package servercommon

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func JsonRequest(method string, url string, data url.Values, result any) error {
	encodedData := data.Encode()
	req, err := http.NewRequest(method, url, strings.NewReader(encodedData))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	if result == nil {
		return nil
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("Request content=%s", string(content)))
	err = json.Unmarshal(content, result)
	return err
}
