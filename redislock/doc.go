// Package redislock provides a Redis distributed locking mechanism.
//
// It borrows some code from https://github.com/bsm/redislock.
// The design is different, and allows simpler mocking.
// It also provides tracing and better error messages.
package redislock
