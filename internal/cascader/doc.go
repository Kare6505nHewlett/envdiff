// Package cascader provides layer-based merging of .env file values.
//
// Layers are applied in order from lowest to highest priority. When the
// same key appears in multiple layers the last layer's value wins. Each
// Result records which layer supplied the final value and whether any
// earlier layer was overridden, making it easy to audit configuration
// drift across environments such as base → staging → production.
//
// Example usage:
//
//	layers := []cascader.Layer{
//		{Name: "base",    Values: baseEnv},
//		{Name: "staging", Values: stagingEnv},
//	}
//	results, err := cascader.Apply(layers, cascader.DefaultOptions())
package cascader
