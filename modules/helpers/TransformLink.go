package helpers

import "net/url"

func TransformLink(link *url.URL) *url.URL {
	transformedLink := link.Scheme + "://" + link.Host + link.Path

	parsedLink, _ := url.Parse(transformedLink)

	return parsedLink
}
