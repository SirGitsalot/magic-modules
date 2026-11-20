package assuredworkloadsv2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestDataSourceAssuredWorkloadsV2WorkloadSchema(t *testing.T) {
	t.Parallel()

	ds := DataSourceAssuredWorkloadsV2Workload()
	if err := ds.InternalValidate(nil, false); err != nil {
		t.Fatalf("data source schema InternalValidate failed: %s", err)
	}
	if _, exists := ds.Schema["resource_config"]; exists {
		t.Errorf("expected %q to be excluded from data source schema", "resource_config")
	}
	if _, exists := ds.Schema["deletion_policy"]; !exists {
		t.Errorf("expected %q to remain in data source schema so resource read can populate virtual fields", "deletion_policy")
	}

	cases := []struct {
		name    string
		raw     map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid name",
			raw: map[string]interface{}{
				"name": "organizations/123/locations/us-east4/workloads/my-workload",
			},
		},
		{
			name: "valid org location workload_id",
			raw: map[string]interface{}{
				"organization": "123",
				"location":     "us-east4",
				"workload_id":  "my-workload",
			},
		},
		{
			name:    "neither specified",
			raw:     map[string]interface{}{},
			wantErr: true,
		},
		{
			name: "invalid name format",
			raw: map[string]interface{}{
				"name": "invalid-name",
			},
			wantErr: true,
		},
		{
			name: "name conflicts with organization",
			raw: map[string]interface{}{
				"name":         "organizations/123/locations/us-east4/workloads/my-workload",
				"organization": "123",
			},
			wantErr: true,
		},
		{
			name: "incomplete components",
			raw: map[string]interface{}{
				"organization": "123",
				"workload_id":  "my-workload",
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			diags := ds.Validate(terraform.NewResourceConfigRaw(tc.raw))
			if tc.wantErr && !diags.HasError() {
				t.Errorf("expected validation error for %#v, got none", tc.raw)
			}
			if !tc.wantErr && diags.HasError() {
				t.Errorf("unexpected validation errors for %#v: %v", tc.raw, diags)
			}
		})
	}
}

func TestDataSourceAssuredWorkloadsV2WorkloadResolveId(t *testing.T) {
	t.Parallel()

	dsSchema := DataSourceAssuredWorkloadsV2Workload().Schema
	wantName := "organizations/123/locations/us-east4/workloads/my-workload"

	t.Run("by name populates regional URL components", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, dsSchema, map[string]interface{}{
			"name": wantName,
		})
		id, err := resolveDataSourceAssuredWorkloadsV2WorkloadId(d)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if id != wantName {
			t.Errorf("id = %q, want %q", id, wantName)
		}
		for k, want := range map[string]string{
			"organization": "123",
			"location":     "us-east4",
			"workload_id":  "my-workload",
			"name":         wantName,
		} {
			if got := d.Get(k).(string); got != want {
				t.Errorf("%s = %q, want %q", k, got, want)
			}
		}
	})

	t.Run("by components populates name", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, dsSchema, map[string]interface{}{
			"organization": "123",
			"location":     "us-east4",
			"workload_id":  "my-workload",
		})
		id, err := resolveDataSourceAssuredWorkloadsV2WorkloadId(d)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if id != wantName {
			t.Errorf("id = %q, want %q", id, wantName)
		}
		if got := d.Get("name").(string); got != wantName {
			t.Errorf("name = %q, want %q", got, wantName)
		}
	})

	t.Run("invalid name returns error", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, dsSchema, map[string]interface{}{
			"name": "not-a-valid-name",
		})
		if _, err := resolveDataSourceAssuredWorkloadsV2WorkloadId(d); err == nil {
			t.Error("expected an error for invalid name")
		}
	})

	t.Run("missing components returns error", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, dsSchema, map[string]interface{}{
			"organization": "123",
		})
		if _, err := resolveDataSourceAssuredWorkloadsV2WorkloadId(d); err == nil {
			t.Error("expected an error for missing components")
		}
	})
}
