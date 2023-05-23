/*
Copyright (c) 2023 SUSE LLC

This program is free software; you can redistribute it and/or
modify it under the terms of version 3 of the GNU General Public License as
published by the Free Software Foundation.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.   See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, contact SUSE LLC.

To contact SUSE about this file by physical or electronic mail,
you may find current contact information at www.suse.com
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

//+genclient
//+genclient:nonNamespaced
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:object:root=true

// ProductUsage is the Schema for the productusages API
type ProductUsage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	BaseProduct      string `json:"base_product"`
        ManagedNodeCount uint32 `json:"managed_node_count"`
        ReportingTime    string `json:"reporting_time"`
}

//+kubebuilder:object:root=true

// ProductUsageList contains a list of ProductUsage
type ProductUsageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProductUsage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProductUsage{}, &ProductUsageList{})
}
