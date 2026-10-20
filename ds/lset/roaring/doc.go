// Package roaring provides a compressed set of uint32 values built on Roaring
// bitmaps, which stay compact and fast for both sparse and dense sets.
//
// See Chambi, Lemire et al., "Better bitmap performance with Roaring bitmaps"
// (2016), https://arxiv.org/abs/1402.6407, and https://roaringbitmap.org/.
package roaring

// == projects.Code.luce.roaring ==
// [ ] Investigate SIMD for container ops
//  Once the basic version works and has benchmarks. Candidates are
//  bitmap-bitmap And/Or/AndNot and array-array intersection. Try the
//  experimental simd/archsimd packages (GOEXPERIMENT=simd) behind a build
//  tag with the scalar version as fallback; assembly (avo) as an
//  alternative. Also measure GOAMD64=v3.
