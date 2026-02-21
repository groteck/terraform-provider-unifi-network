package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringPtr returns a pointer to the string value if it's not null or unknown.
func StringPtr(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	s := v.ValueString()
	return &s
}

// StringValue returns a types.String from a string pointer.
func StringValue(v *string) types.String {
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(*v)
}

// StringToValue returns a types.String from a string. Returns null if empty.
func StringToValue(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

// Int64Ptr returns a pointer to the int64 value if it's not null or unknown.
func Int64Ptr(v types.Int64) *int {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	i := int(v.ValueInt64())
	return &i
}

// Int64Value returns a types.Int64 from an int pointer.
func Int64Value(v *int) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*v))
}

// BoolPtr returns a pointer to the bool value if it's not null or unknown.
func BoolPtr(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	b := v.ValueBool()
	return &b
}

// BoolValue returns a types.Bool from a bool pointer.
func BoolValue(v *bool) types.Bool {
	if v == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*v)
}

// StringOrEmpty returns the string value or empty if null/unknown.
func StringOrEmpty(v types.String) string {
	if v.IsNull() || v.IsUnknown() {
		return ""
	}
	return v.ValueString()
}
