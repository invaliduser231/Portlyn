package domain

import "testing"

func TestNormalizeSubdomain(t *testing.T) {
	value, err := NormalizeSubdomain(" Pangolin.App ")
	if err != nil {
		t.Fatalf("normalize subdomain: %v", err)
	}
	if value != "pangolin.app" {
		t.Fatalf("unexpected normalized subdomain %q", value)
	}
}

func TestServiceHostname(t *testing.T) {
	if got := ServiceHostname("schnittert.cloud", "pangolin"); got != "pangolin.schnittert.cloud" {
		t.Fatalf("unexpected hostname %q", got)
	}
	if got := ServiceHostname("schnittert.cloud", ""); got != "schnittert.cloud" {
		t.Fatalf("unexpected root hostname %q", got)
	}
}
