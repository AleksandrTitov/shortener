package url

import "net/url"

func ValidateURL(u string) error {
	_, err := url.ParseRequestURI(u)
	return err
}
