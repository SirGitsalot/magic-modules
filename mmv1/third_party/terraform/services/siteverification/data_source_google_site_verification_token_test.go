// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package siteverification_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSiteVerificationToken_siteverificationTokenSite(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "siteverify", serviceAccount)

	context := map[string]interface{}{
		"site":    "https://www.example.com",
		"account": targetServiceAccountEmail,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSiteVerificationToken_siteverificationTokenSite(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_site_verification_token.site_meta", "token"),
					resource.TestCheckResourceAttr("data.google_site_verification_token.site_meta", "type", "SITE"),
					resource.TestCheckResourceAttr("data.google_site_verification_token.site_meta", "identifier", context["site"].(string)),
					resource.TestCheckResourceAttr("data.google_site_verification_token.site_meta", "verification_method", "META"),
				),
			},
		},
	})
}

func testAccSiteVerificationToken_siteverificationTokenSite(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_service_account_access_token" "impersonated" {
  target_service_account = "%{account}"
  scopes = [
    "https://www.googleapis.com/auth/siteverification",
    "https://www.googleapis.com/auth/siteverification.verify_only",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/userinfo.email",
  ]
  lifetime = "300s"
}

provider "google" {
  alias                 = "impersonated"
  user_project_override = true
  access_token          = data.google_service_account_access_token.impersonated.access_token
}

data "google_site_verification_token" "site_meta" {
  provider            = google.impersonated
  type                = "SITE"
  identifier          = "%{site}"
  verification_method = "META"
}
`, context)
}

func TestAccSiteVerificationToken_siteverificationTokenDomain(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "siteverify", serviceAccount)

	context := map[string]interface{}{
		"domain":  "www.example.com",
		"account": targetServiceAccountEmail,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSiteVerificationToken_siteverificationTokenDomain(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_site_verification_token.dns_text", "token"),
					resource.TestCheckResourceAttr("data.google_site_verification_token.dns_text", "type", "INET_DOMAIN"),
					resource.TestCheckResourceAttr("data.google_site_verification_token.dns_text", "identifier", context["domain"].(string)),
					resource.TestCheckResourceAttr("data.google_site_verification_token.dns_text", "verification_method", "DNS_TXT"),
				),
			},
		},
	})
}

func testAccSiteVerificationToken_siteverificationTokenDomain(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_service_account_access_token" "impersonated" {
  target_service_account = "%{account}"
  scopes = [
    "https://www.googleapis.com/auth/siteverification",
    "https://www.googleapis.com/auth/siteverification.verify_only",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/userinfo.email",
  ]
  lifetime = "300s"
}

provider "google" {
  alias                 = "impersonated"
  user_project_override = true
  access_token          = data.google_service_account_access_token.impersonated.access_token
}

data "google_site_verification_token" "dns_text" {
  provider            = google.impersonated
  type                = "INET_DOMAIN"
  identifier          = "%{domain}"
  verification_method = "DNS_TXT"
}
`, context)
}
