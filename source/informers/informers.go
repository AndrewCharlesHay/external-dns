/*
Copyright 2025 The Kubernetes Authors.

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

package informers

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type informerFactory interface {
	WaitForCacheSync(stopCh <-chan struct{}) map[reflect.Type]bool
}

type dynamicInformerFactory interface {
	WaitForCacheSync(stopCh <-chan struct{}) map[schema.GroupVersionResource]bool
}

// DefaultCacheSyncTimeout is the default timeout for waiting for informer
// caches to sync. A timeout is always enforced to prevent indefinite hangs
// caused by misconfiguration (e.g. missing RBAC or CRD not installed).
const DefaultCacheSyncTimeout = 60 * time.Second

// WaitForCacheSync waits for all informers in the factory to sync their caches.
// If timeout is <= 0, the default timeout (60s) is used.
// Returns an error if any informer fails to sync within the timeout.
func WaitForCacheSync(ctx context.Context, factory informerFactory, timeout time.Duration) error {
	return waitForCacheSync(ctx, factory.WaitForCacheSync, timeout)
}

// WaitForDynamicCacheSync waits for all dynamic informers in the factory to sync their caches.
// If timeout is <= 0, the default timeout (60s) is used.
// Returns an error if any informer fails to sync within the timeout.
func WaitForDynamicCacheSync(ctx context.Context, factory dynamicInformerFactory, timeout time.Duration) error {
	return waitForCacheSync(ctx, factory.WaitForCacheSync, timeout)
}

// waitForCacheSync waits for informer caches to sync within the given timeout.
// If timeout is <= 0, the default (60s) is used to prevent indefinite hangs.
// Returns an error if any cache fails to sync.
func waitForCacheSync[K comparable](ctx context.Context, waitFunc func(<-chan struct{}) map[K]bool, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = DefaultCacheSyncTimeout
	}
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()
	for typ, done := range waitFunc(ctx.Done()) {
		if !done {
			if ctx.Err() != nil {
				return fmt.Errorf("failed to sync %v: %w", typ, ctx.Err())
			}
			return fmt.Errorf("failed to sync %v", typ)
		}
	}
	return nil
}
