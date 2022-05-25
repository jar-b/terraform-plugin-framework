package fromproto6_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fromproto6"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func TestUpgradeResourceStateRequest(t *testing.T) {
	t.Parallel()

	testFwSchema := &tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"test_attribute": {
				Required: true,
				Type:     types.StringType,
			},
		},
	}

	testCases := map[string]struct {
		input               *tfprotov6.UpgradeResourceStateRequest
		resourceSchema      *tfsdk.Schema
		resourceType        tfsdk.ResourceType
		expected            *fwserver.UpgradeResourceStateRequest
		expectedDiagnostics diag.Diagnostics
	}{
		"nil": {
			input:    nil,
			expected: nil,
		},
		"rawstate": {
			input: &tfprotov6.UpgradeResourceStateRequest{
				RawState: testNewRawState(t, map[string]interface{}{
					"test_attribute": "test-value",
				}),
			},
			resourceSchema: testFwSchema,
			expected: &fwserver.UpgradeResourceStateRequest{
				RawState: testNewRawState(t, map[string]interface{}{
					"test_attribute": "test-value",
				}),
				ResourceSchema: *testFwSchema,
			},
		},
		"resourceschema": {
			input:          &tfprotov6.UpgradeResourceStateRequest{},
			resourceSchema: testFwSchema,
			expected: &fwserver.UpgradeResourceStateRequest{
				ResourceSchema: *testFwSchema,
			},
		},
		"resourceschema-missing": {
			input:    &tfprotov6.UpgradeResourceStateRequest{},
			expected: nil,
			expectedDiagnostics: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Create Empty State",
					"An unexpected error was encountered when creating the empty state. "+
						"This is always an issue in the Terraform Provider SDK used to implement the provider and should be reported to the provider developers.\n\n"+
						"Please report this to the provider developer:\n\n"+
						"Missing schema.",
				),
			},
		},
		"version": {
			input: &tfprotov6.UpgradeResourceStateRequest{
				Version: 123,
			},
			resourceSchema: testFwSchema,
			expected: &fwserver.UpgradeResourceStateRequest{
				ResourceSchema: *testFwSchema,
				Version:        123,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, diags := fromproto6.UpgradeResourceStateRequest(context.Background(), testCase.input, testCase.resourceType, testCase.resourceSchema)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}