// Package go_persistent_ds implements several fully persistent data structures.
//
// Available data structures are:
//   - Map
//   - Slice
//   - DoubleLinkedList
//
// All structures are base on FatNodes.
//
// Note that every structure can perform total of 2^65-1 modifications, and will panic on attempt to modify it for 2^65 time.
// If you need to continue editing structure, the good idea is to use appropriate method to dump structure for special version.
package go_persistent_ds
