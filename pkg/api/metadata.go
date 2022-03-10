/*
 * Copyright (c) 2017, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"fmt"
	"net/http"
	"reflect"
	"sort"

	"github.com/go-chi/chi/v5"
	yaml "gopkg.in/yaml.v2"

	"github.com/megaease/easegress/pkg/filters"
	"github.com/megaease/easegress/pkg/v"
)

const (
	// MetadataPrefix is the metadata prefix.
	MetadataPrefix = "/metadata"

	// ObjectMetadataPrefix is the object metadata prefix.
	ObjectMetadataPrefix = "/metadata/objects"

	// FilterMetaPrefix is the filter of Pipeline metadata prefix.
	FilterMetaPrefix = "/metadata/objects/pipeline/filters"
)

func (s *Server) metadataAPIEntries() []*Entry {
	return []*Entry{
		{
			Path:    FilterMetaPrefix,
			Method:  "GET",
			Handler: s.listFilters,
		},
		{
			Path:    FilterMetaPrefix + "/{kind}" + "/description",
			Method:  "GET",
			Handler: s.getFilterDescription,
		},
		{
			Path:    FilterMetaPrefix + "/{kind}" + "/schema",
			Method:  "GET",
			Handler: s.getFilterSchema,
		},
		{
			Path:    FilterMetaPrefix + "/{kind}" + "/results",
			Method:  "GET",
			Handler: s.getFilterResults,
		},
	}
}

func (s *Server) listFilters(w http.ResponseWriter, r *http.Request) {
	var kinds []string
	filters.WalkKind(func(k *filters.Kind) bool {
		kinds = append(kinds, k.Name)
		return true
	})
	sort.Strings(kinds)

	buff, err := yaml.Marshal(kinds)
	if err != nil {
		panic(fmt.Errorf("marshal %#v to yaml failed: %v", kinds, err))
	}

	w.Header().Set("Content-Type", "text/vnd.yaml")
	w.Write(buff)
}

func (s *Server) getFilterDescription(w http.ResponseWriter, r *http.Request) {
	kind := chi.URLParam(r, "kind")

	k := filters.GetKind(kind)
	if k == nil {
		HandleAPIError(w, r, http.StatusNotFound, fmt.Errorf("not found"))
		return
	}
	w.Write([]byte(k.Description))
}

func (s *Server) getFilterSchema(w http.ResponseWriter, r *http.Request) {
	kind := chi.URLParam(r, "kind")

	k := filters.GetKind(kind)
	if k == nil {
		HandleAPIError(w, r, http.StatusNotFound, fmt.Errorf("not found"))
		return
	}
	specType := reflect.TypeOf(k.DefaultSpec())

	buff, err := v.GetSchemaInYAML(specType)
	if err != nil {
		panic(fmt.Errorf("get schema for %v failed: %v", kind, err))
	}

	w.Header().Set("Content-Type", "text/vnd.yaml")
	w.Write(buff)
}

func (s *Server) getFilterResults(w http.ResponseWriter, r *http.Request) {
	kind := chi.URLParam(r, "kind")

	k := filters.GetKind(kind)
	if k == nil {
		HandleAPIError(w, r, http.StatusNotFound, fmt.Errorf("not found"))
		return
	}

	results := append([]string{}, k.Results...)
	sort.Strings(results)

	buff, err := yaml.Marshal(results)
	if err != nil {
		panic(fmt.Errorf("marshal %#v to yaml failed: %v", results, err))
	}

	w.Header().Set("Content-Type", "text/vnd.yaml")
	w.Write(buff)
}
