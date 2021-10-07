/*
Copyright 2021.

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

package utils

import (
	"math"
	"sync"
	"time"
)

type DefaultRateLimiter struct {
	requeuesLock sync.Mutex
	requeues     map[interface{}]int
	modifier     map[interface{}]int

	baseDelay time.Duration
	maxDelay  time.Duration
}

func NewDefaultRateLimiter(baseDelay, maxDelay time.Duration) *DefaultRateLimiter {
	return &DefaultRateLimiter{
		baseDelay: baseDelay,
		maxDelay:  maxDelay,
		requeues:  map[interface{}]int{},
		modifier:  map[interface{}]int{},
	}
}

func (r *DefaultRateLimiter) When(item interface{}) time.Duration {
	r.requeuesLock.Lock()
	defer r.requeuesLock.Unlock()

	exp := r.modifier[item]
	r.requeues[item]++

	if r.requeues[item]%16 == 0 {
		r.modifier[item]++
	}

	// The backoff is capped such that 'calculated' value never overflows.
	backoff := float64(r.baseDelay.Nanoseconds()) * math.Pow(2, float64(exp))
	if backoff > math.MaxInt64 {
		return r.maxDelay
	}

	calculated := time.Duration(backoff)
	if calculated > r.maxDelay {
		return r.maxDelay
	}

	return calculated
}

func (r *DefaultRateLimiter) NumRequeues(item interface{}) int {
	r.requeuesLock.Lock()
	defer r.requeuesLock.Unlock()

	return r.requeues[item]
}

func (r *DefaultRateLimiter) Forget(item interface{}) {
	r.requeuesLock.Lock()
	defer r.requeuesLock.Unlock()

	delete(r.requeues, item)
}
