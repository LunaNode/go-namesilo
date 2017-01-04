// Copyright (c) 2017 LunaNode Hosting Inc. All right reserved.
// Use of this source code is governed by the MIT License. See LICENSE file.

package namesilo

type Domain struct {
	Domain string
	client *Client
}

type listDomainsResponse struct {
	Domains []string `xml:"domains>domain"`
}

func (client *Client) ListDomains() ([]Domain, error) {
	var response listDomainsResponse
	if err := client.request("listDomains", nil, &response); err != nil {
		return nil, err
	}
	var domains []Domain
	for _, item := range response.Domains {
		domains = append(domains, Domain{
			Domain: item,
			client: client,
		})
	}
	return domains, nil
}
