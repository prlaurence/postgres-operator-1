/*
Copyright 2020 Crunchy Data Solutions, Inc.
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
	v1 "github.com/crunchydata/postgres-operator/apis/crunchydata.com/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// PgpolicyLister helps list Pgpolicies.
type PgpolicyLister interface {
	// List lists all Pgpolicies in the indexer.
	List(selector labels.Selector) (ret []*v1.Pgpolicy, err error)
	// Pgpolicies returns an object that can list and get Pgpolicies.
	Pgpolicies(namespace string) PgpolicyNamespaceLister
	PgpolicyListerExpansion
}

// pgpolicyLister implements the PgpolicyLister interface.
type pgpolicyLister struct {
	indexer cache.Indexer
}

// NewPgpolicyLister returns a new PgpolicyLister.
func NewPgpolicyLister(indexer cache.Indexer) PgpolicyLister {
	return &pgpolicyLister{indexer: indexer}
}

// List lists all Pgpolicies in the indexer.
func (s *pgpolicyLister) List(selector labels.Selector) (ret []*v1.Pgpolicy, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Pgpolicy))
	})
	return ret, err
}

// Pgpolicies returns an object that can list and get Pgpolicies.
func (s *pgpolicyLister) Pgpolicies(namespace string) PgpolicyNamespaceLister {
	return pgpolicyNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// PgpolicyNamespaceLister helps list and get Pgpolicies.
type PgpolicyNamespaceLister interface {
	// List lists all Pgpolicies in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Pgpolicy, err error)
	// Get retrieves the Pgpolicy from the indexer for a given namespace and name.
	Get(name string) (*v1.Pgpolicy, error)
	PgpolicyNamespaceListerExpansion
}

// pgpolicyNamespaceLister implements the PgpolicyNamespaceLister
// interface.
type pgpolicyNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Pgpolicies in the indexer for a given namespace.
func (s pgpolicyNamespaceLister) List(selector labels.Selector) (ret []*v1.Pgpolicy, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Pgpolicy))
	})
	return ret, err
}

// Get retrieves the Pgpolicy from the indexer for a given namespace and name.
func (s pgpolicyNamespaceLister) Get(name string) (*v1.Pgpolicy, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("pgpolicy"), name)
	}
	return obj.(*v1.Pgpolicy), nil
}
