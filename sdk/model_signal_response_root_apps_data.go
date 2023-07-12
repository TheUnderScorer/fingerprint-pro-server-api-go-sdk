/*
 * Fingerprint Pro Server API
 *
 * Fingerprint Pro Server API allows you to get information about visitors and about individual events in a server environment. It can be used for data exports, decision-making, and data analysis scenarios. Server API is intended for server-side usage, it's not intended to be used from the client side, whether it's a browser or a mobile device.
 *
 * API version: 3
 * Contact: support@fingerprint.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package sdk

type SignalResponseRootAppsData struct {
	// Android specific root management apps detection. There are 2 values: • `true` - Root Management Apps detected (e.g. Magisk) • `false` - No Root Management Apps detected or the client is not Android.
	Result bool `json:"result,omitempty"`
}
