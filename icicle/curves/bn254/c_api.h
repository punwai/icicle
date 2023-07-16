
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

#include <stdbool.h>
#include <cuda.h>
// c_api.h

#ifdef __cplusplus
extern "C" {
#endif

typedef struct BN254_projective_t BN254_projective_t;
typedef struct  BN254_g2_projective_t BN254_g2_projective_t;

bool eq_bn254(BN254_projective_t *point1, BN254_projective_t *point2);
bool eq_g2_bn254(BN254_g2_projective_t *point1, BN254_g2_projective_t *point2);

#ifdef __cplusplus
}
#endif
