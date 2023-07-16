
	// Copyright 2023 Ingonyama
	//
	// Licensed under the Apache License, Version 2.0 (the "License");
	// you may not use this file except in compliance with the License.
	// You may obtain a copy of the License at
	//
	//     http://www.apache.org/licenses/LICENSE-2.0
	//
	// Unless required by applicable law or agreed to in writing, software
	// distributed under the License is distributed on an "AS IS" BASIS,
	// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	// See the License for the specific language governing permissions and
	// limitations under the License.
	
// Code generated by Ingonyama DO NOT EDIT

package bls12381

import (
	"unsafe"

	"encoding/binary"
	"fmt"

	


	"github.com/consensys/gnark-crypto/ecc/bls12-381"



	

	"github.com/consensys/gnark-crypto/ecc/bls12-381/fp"



	


	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"



)

// #cgo CFLAGS: -I${SRCDIR}/icicle/curves/bls12381/
// #cgo LDFLAGS: -L${SRCDIR}/../../ -lbn12_381
// #include "c_api.h"
// #include "ve_mod_mult.h"
import "C"

const SCALAR_SIZE = 8
const BASE_SIZE = 12

type ScalarField struct {
	s [SCALAR_SIZE]uint32
}

type BaseField struct {
	s [BASE_SIZE]uint32
}

type Field interface {
	toGnarkFr() *fr.Element
}

/*
 * Common Constrctors
 */

func NewFieldZero[T BaseField | ScalarField]() *T {
	var field T

	return &field
}

func NewFieldFromFrGnark[T BaseField | ScalarField](element fr.Element) *T {
	s := ConvertUint64ArrToUint32Arr(element.Bits()) // get non-montgomry

	return &T{s}
}

func NewFieldFromFpGnark[T BaseField | ScalarField](element fp.Element) *T {
	s := ConvertUint64ArrToUint32Arr(element.Bits()) // get non-montgomry

	return &T{s}
}

/*
 * BaseField Constrctors
 */

func NewBaseFieldOne() *BaseField {
	var s [BASE_SIZE]uint32

	s[0] = 1

	return &BaseField{s}
}

func BaseFieldFromLimbs(limbs [BASE_SIZE]uint32) *BaseField {
	bf := NewFieldZero[BaseField]()
	copy(bf.s[:], limbs[:])

	return bf
}

/*
 * BaseField methods
 */

func (f *BaseField) limbs() [BASE_SIZE]uint32 {
	return f.s
}

func (f *BaseField) toBytesLe() []byte {
	bytes := make([]byte, len(f.s)*4)
	for i, v := range f.s {
		binary.LittleEndian.PutUint32(bytes[i*4:], v)
	}

	return bytes
}

func (f *BaseField) toGnarkFr() *fr.Element {
	fb := f.toBytesLe()
	var b32 [32]byte
	copy(b32[:], fb[:32])

	v, e := fr.LittleEndian.Element(&b32)

	if e != nil {
		panic(fmt.Sprintf("unable to create convert point %v got error %v", f, e))
	}

	return &v
}

func (f *BaseField) toGnarkFp() *fp.Element {
	fb := f.toBytesLe()
	var b32 [32]byte
	copy(b32[:], fb[:32])

	v, e := fp.LittleEndian.Element(&b32)

	if e != nil {
		panic(fmt.Sprintf("unable to create convert point %v got error %v", f, e))
	}

	return &v
}

/*
 * ScalarField methods
 */

func NewScalarFieldOne() *ScalarField {
	var s [SCALAR_SIZE]uint32

	s[0] = 1

	return &ScalarField{s}
}

/*
 * ScalarField methods
 */

func (f *ScalarField) limbs() [SCALAR_SIZE]uint32 {
	return f.s
}

func (f *ScalarField) toBytesLe() []byte {
	bytes := make([]byte, len(f.s)*4)
	for i, v := range f.s {
		binary.LittleEndian.PutUint32(bytes[i*4:], v)
	}

	return bytes
}

func (f ScalarField) toGnarkFr() *fr.Element {
	fb := f.toBytesLe()
	var b32 [32]byte
	copy(b32[:], fb[:32])

	v, e := fr.LittleEndian.Element(&b32)

	if e != nil {
		panic(fmt.Sprintf("unable to create convert point %v got error %v", f, e))
	}

	return &v
}

func (f *ScalarField) toGnarkFp() *fp.Element {
	fb := f.toBytesLe()
	var b32 [32]byte
	copy(b32[:], fb[:32])

	v, e := fp.LittleEndian.Element(&b32)

	if e != nil {
		panic(fmt.Sprintf("unable to create convert point %v got error %v", f, e))
	}

	return &v
}

/*
 * PointBLS12381
 */

type PointBLS12381 struct {
	x, y, z BaseField
}

func NewPointBLS12381Zero() *PointBLS12381 {
	return &PointBLS12381{
		x: *NewFieldZero[BaseField](),
		y: *NewBaseFieldOne(),
		z: *NewFieldZero[BaseField](),
	}
}

func (p *PointBLS12381) eq(pCompare *PointBLS12381) bool {
	// Cast *PointBLS12381 to *C.BLS12381_projective_t
	// The unsafe.Pointer cast is necessary because Go doesn't allow direct casts
	// between different pointer types.
	// It's your responsibility to ensure that the types are compatible.
	pC := (*C.BLS12381_projective_t)(unsafe.Pointer(p))
	pCompareC := (*C.BLS12381_projective_t)(unsafe.Pointer(pCompare))

	// Call the C function
	// The C function doesn't keep any references to the data,
	// so it's fine if the Go garbage collector moves or deletes the data later.
	return bool(C.eq_bls12381(pC, pCompareC))
}

func (p *PointBLS12381) strip_z() *PointAffineNoInfinityBLS12381 {
	return &PointAffineNoInfinityBLS12381{
		x: p.x,
		y: p.y,
	}
}

func (p *PointBLS12381) toGnarkAffine() *bls12381.G1Affine {
	px := p.x.toGnarkFp()
	py := p.y.toGnarkFp()
	pz := p.z.toGnarkFp()

	zInv := new(fp.Element)
	x := new(fp.Element)
	y := new(fp.Element)

	zInv.Inverse(pz)

	x.Mul(px, zInv)
	y.Mul(py, zInv)

	return &bls12381.G1Affine{X: *x, Y: *y}
}

func (p *PointBLS12381) ToGnarkJac() *bls12381.G1Jac {
	var p1 bls12381.G1Jac
	p1.FromAffine(p.toGnarkAffine())

	return &p1
}

func PointBLS12381FromG1AffineGnark(gnark *bls12381.G1Affine) *PointBLS12381 {
	point := PointBLS12381{
		x: *NewFieldFromFpGnark[BaseField](gnark.X),
		y: *NewFieldFromFpGnark[BaseField](gnark.Y),
		z: *NewBaseFieldOne(),
	}

	return &point
}

// converts jac fromat to projective
func PointBLS12381FromJacGnark(gnark *bls12381.G1Jac) *PointBLS12381 {
	var pointAffine bls12381.G1Affine
	pointAffine.FromJacobian(gnark)

	point := PointBLS12381{
		x: *NewFieldFromFpGnark[BaseField](pointAffine.X),
		y: *NewFieldFromFpGnark[BaseField](pointAffine.Y),
		z: *NewBaseFieldOne(),
	}

	return &point
}

func PointBLS12381fromLimbs(x, y, z *[]uint32) *PointBLS12381 {
	return &PointBLS12381{
		x: *BaseFieldFromLimbs(getFixedLimbs(x)),
		y: *BaseFieldFromLimbs(getFixedLimbs(y)),
		z: *BaseFieldFromLimbs(getFixedLimbs(z)),
	}
}

/*
 * PointAffineNoInfinityBLS12381
 */

type PointAffineNoInfinityBLS12381 struct {
	x, y BaseField
}

func NewPointAffineNoInfinityBLS12381Zero() *PointAffineNoInfinityBLS12381 {
	return &PointAffineNoInfinityBLS12381{
		x: *NewFieldZero[BaseField](),
		y: *NewFieldZero[BaseField](),
	}
}

func (p *PointAffineNoInfinityBLS12381) toProjective() *PointBLS12381 {
	return &PointBLS12381{
		x: p.x,
		y: p.y,
		z: *NewBaseFieldOne(),
	}
}

func (p *PointAffineNoInfinityBLS12381) toGnarkAffine() *bls12381.G1Affine {
	return p.toProjective().toGnarkAffine()
}

func PointAffineNoInfinityBLS12381FromLimbs(x, y *[]uint32) *PointAffineNoInfinityBLS12381 {
	return &PointAffineNoInfinityBLS12381{
		x: *BaseFieldFromLimbs(getFixedLimbs(x)),
		y: *BaseFieldFromLimbs(getFixedLimbs(y)),
	}
}

/*
 * Multiplication
 */

func MultiplyVec(a []PointBLS12381, b []ScalarField, deviceID int) {
	if len(a) != len(b) {
		panic("a and b have different lengths")
	}

	pointsC := (*C.BLS12381_projective_t)(unsafe.Pointer(&a[0]))
	scalarsC := (*C.BLS12381_scalar_t)(unsafe.Pointer(&b[0]))
	deviceIdC := C.size_t(deviceID)
	nElementsC := C.size_t(len(a))

	C.vec_mod_mult_point_bls12381(pointsC, scalarsC, nElementsC, deviceIdC)
}

func MultiplyScalar(a []ScalarField, b []ScalarField, deviceID int) {
	if len(a) != len(b) {
		panic("a and b have different lengths")
	}

	aC := (*C.BLS12381_scalar_t)(unsafe.Pointer(&a[0]))
	bC := (*C.BLS12381_scalar_t)(unsafe.Pointer(&b[0]))
	deviceIdC := C.size_t(deviceID)
	nElementsC := C.size_t(len(a))

	C.vec_mod_mult_scalar_bls12381(aC, bC, nElementsC, deviceIdC)
}

// Multiply a matrix by a scalar:
//
//	`a` - flattenned matrix;
//	`b` - vector to multiply `a` by;
func MultiplyMatrix(a []ScalarField, b []ScalarField, deviceID int) {
	c := make([]ScalarField, len(b))
	for i := range c {
		c[i] = *NewFieldZero[ScalarField]()
	}

	aC := (*C.BLS12381_scalar_t)(unsafe.Pointer(&a[0]))
	bC := (*C.BLS12381_scalar_t)(unsafe.Pointer(&b[0]))
	cC := (*C.BLS12381_scalar_t)(unsafe.Pointer(&c[0]))
	deviceIdC := C.size_t(deviceID)
	nElementsC := C.size_t(len(a))

	C.matrix_vec_mod_mult_bls12381(aC, bC, cC, nElementsC, deviceIdC)
}

/*
 * Utils
 */

func getFixedLimbs(slice *[]uint32) [BASE_SIZE]uint32 {
	if len(*slice) <= BASE_SIZE {
		limbs := [BASE_SIZE]uint32{}
		copy(limbs[:len(*slice)], *slice)
		return limbs
	}

	panic("slice has too many elements")
}

func BatchConvertFromFrGnark[T BaseField | ScalarField](elements []fr.Element) []T {
	var newElements []T
	for _, e := range elements {
		converted := NewFieldFromFrGnark[T](e)
		newElements = append(newElements, *converted)
	}

	return newElements
}

func BatchConvertFromFrGnarkThreaded[T BaseField | ScalarField](elements []fr.Element, routines int) []T {
	var newElements []T

	if routines > 1 {
		channels := make([]chan []T, routines)
		for i := 0; i < routines; i ++ {
			channels[i] = make(chan []T, 1)
		} 

		convert := func(elements []fr.Element, chanIndex int) {
			var convertedElements []T
			for _, e := range elements {
				converted := NewFieldFromFrGnark[T](e)
				convertedElements = append(convertedElements, *converted)
			}

			channels[chanIndex] <- convertedElements
		}

		batchLen := len(elements)/routines
		for i := 0; i < routines; i ++ {
			elemsToConv := elements[batchLen*i:batchLen*(i+1)]
			go convert(elemsToConv, i)
		}

		for i := 0; i < routines; i ++ {
			newElements = append(newElements, <-channels[i]...)
		}
	} else {
		for _, e := range elements {
			converted := NewFieldFromFrGnark[T](e)
			newElements = append(newElements, *converted)
		}
	}

	return newElements
}

func BatchConvertToFrGnark[T Field](elements []T) []fr.Element {
	var newElements []fr.Element
	for _, e := range elements {
		converted := e.toGnarkFr()
		newElements = append(newElements, *converted)
	}

	return newElements
}

func BatchConvertToFrGnarkThreaded[T Field](elements []T, routines int) []fr.Element {
	var newElements []fr.Element

	if routines > 1 {
		channels := make([]chan []fr.Element, routines)
		for i := 0; i < routines; i ++ {
			channels[i] = make(chan []fr.Element, 1)
		} 

		convert := func(elements []T, chanIndex int) {
			var convertedElements []fr.Element
			for _, e := range elements {
				converted := e.toGnarkFr()
				convertedElements = append(convertedElements, *converted)
			}

			channels[chanIndex] <- convertedElements
		}

		batchLen := len(elements)/routines
		for i := 0; i < routines; i ++ {
			elemsToConv := elements[batchLen*i:batchLen*(i+1)]
			go convert(elemsToConv, i)
		}

		for i := 0; i < routines; i ++ {
			newElements = append(newElements, <-channels[i]...)
		}
	} else {
		for _, e := range elements {
			converted := e.toGnarkFr()
			newElements = append(newElements, *converted)
		}
	}

	return newElements
}

func BatchConvertFromG1Affine(elements []bls12381.G1Affine) []PointAffineNoInfinityBLS12381 {
	var newElements []PointAffineNoInfinityBLS12381
	for _, e := range elements {
		newElement := PointBLS12381FromG1AffineGnark(&e).strip_z()
		newElements = append(newElements, *newElement)
	}
	return newElements
}
