// Copyright (c) 2017 LunaNode Hosting Inc. All right reserved.
// Use of this source code is governed by the MIT License. See LICENSE file.

package namesilo

import (
	"testing"
)

func TestListDomains(t *testing.T) {
	// XML from https://www.namesilo.com/api_reference.php#listDomains
	client := GetTestClient(t, "/listDomains", `<namesilo>
    <request>
        <operation>listDomains</operation>
        <ip>55.555.55.55</ip>
    </request>
    <reply>
        <code>300</code>
        <detail>success</detail>
        <domains>
            <domain>namesilo.com</domain>
            <domain>namesilo.net</domain>
            <domain>namesilo.biz</domain>
            <domain>namesilo.info</domain>
            <domain>namesilo.mobi</domain>
            <domain>namesilo.org</domain>
        </domains>
    </reply>
</namesilo>`)
	domains, err := client.ListDomains()
	if err != nil {
		t.Fatalf("unexpected ListDomains error: %v", err)
	}
	if len(domains) != 6 || domains[1].Domain != "namesilo.net" {
		t.Fatalf("response mismatch (%v)", domains)
	}
}
