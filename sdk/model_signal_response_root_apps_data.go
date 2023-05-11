/*
 * Fingerprint Pro Server API
 *
 * Fingerprint Pro Server API allows you to get information about visitors and about individual events in a server environment. This API can be used for data exports, decision-making, and data analysis scenarios.
 *
 * API version: 3
 * Contact: support@fingerprint.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package sdk

type SignalResponseRootAppsData struct {
	// Android specific root management apps detection. There are 2 values: • `true` - Root Management Apps detected (e.g. Magisk) • `false` - No Root Management Apps detected Available only for events from Android client. The field will not be present for a browser or iOS event.
	Result bool `json:"result,omitempty"`
}
