/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/rancher/csp-rancher-usage-operator/api/csp/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// RancherUsageRecordLister helps list RancherUsageRecords.
// All objects returned here must be treated as read-only.
type RancherUsageRecordLister interface {
	// List lists all RancherUsageRecords in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.RancherUsageRecord, err error)
	// RancherUsageRecords returns an object that can list and get RancherUsageRecords.
	RancherUsageRecords(namespace string) RancherUsageRecordNamespaceLister
	RancherUsageRecordListerExpansion
}

// rancherUsageRecordLister implements the RancherUsageRecordLister interface.
type rancherUsageRecordLister struct {
	indexer cache.Indexer
}

// NewRancherUsageRecordLister returns a new RancherUsageRecordLister.
func NewRancherUsageRecordLister(indexer cache.Indexer) RancherUsageRecordLister {
	return &rancherUsageRecordLister{indexer: indexer}
}

// List lists all RancherUsageRecords in the indexer.
func (s *rancherUsageRecordLister) List(selector labels.Selector) (ret []*v1.RancherUsageRecord, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.RancherUsageRecord))
	})
	return ret, err
}

// RancherUsageRecords returns an object that can list and get RancherUsageRecords.
func (s *rancherUsageRecordLister) RancherUsageRecords(namespace string) RancherUsageRecordNamespaceLister {
	return rancherUsageRecordNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// RancherUsageRecordNamespaceLister helps list and get RancherUsageRecords.
// All objects returned here must be treated as read-only.
type RancherUsageRecordNamespaceLister interface {
	// List lists all RancherUsageRecords in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.RancherUsageRecord, err error)
	// Get retrieves the RancherUsageRecord from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.RancherUsageRecord, error)
	RancherUsageRecordNamespaceListerExpansion
}

// rancherUsageRecordNamespaceLister implements the RancherUsageRecordNamespaceLister
// interface.
type rancherUsageRecordNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all RancherUsageRecords in the indexer for a given namespace.
func (s rancherUsageRecordNamespaceLister) List(selector labels.Selector) (ret []*v1.RancherUsageRecord, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.RancherUsageRecord))
	})
	return ret, err
}

// Get retrieves the RancherUsageRecord from the indexer for a given namespace and name.
func (s rancherUsageRecordNamespaceLister) Get(name string) (*v1.RancherUsageRecord, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("rancherusagerecord"), name)
	}
	return obj.(*v1.RancherUsageRecord), nil
}