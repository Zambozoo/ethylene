#ifndef clox_value_h
#define clox_value_h

#include<stdint.h>

#include "common.h"
#include "floats.h"

typedef uint8_t bit8;
typedef uint16_t bit16;
typedef uint32_t bit32;
typedef uint64_t bit64;

typedef enum {
  VAL_BOOL,
  VAL_NULL, 
  VAL_FLT,
  VAL_INT,
  VAL_FIX,
  VAL_BIT
} ValueType;

typedef struct {
  ValueType type;
  union {
    bool boolean;

    int8_t int8;
    int16_t int16;
    int32_t int32;
    int64_t int64;

    bit8 _bit8;
    bit16 _bit16;
    bit32 _bit32;
    bit64 _bit64;

    float8 flt8;
    float16 flt16;
    float32 flt32;
    double flt64;
  } as; 
} Value;

#define IS_BOOL(value)    ((value).type == VAL_BOOL)
#define IS_NULL(value)     ((value).type == VAL_NIL)
#define IS_FLT(value)  ((value).type == VAL_FLT)
#define IS_INT(value)  ((value).type == VAL_INT)
#define IS_BIT(value)  ((value).type == VAL_BIT)


#define AS_BOOL(value)    ((value).as.boolean)
#define AS_FLT(value)  ((value).as.flt64)
#define AS_INT(value)  ((value).as.int64)
#define AS_BIT(value)  ((value).as._bit64)

#define BOOL_VAL(value) ((Value){VAL_BOOL, {.boolean = value}})
#define NULL_VAL         ((Value){VAL_NULL, {.int64 = 0}})
#define FLT_VAL(value)  ((Value){VAL_FLT, {.flt64 = value}})
#define INT_VAL(value)  ((Value){VAL_INT, {.int64 = value}})
#define BIT_VAL(value)  ((Value){VAL_BIT, {._bit64 = value}})

typedef struct {
  int capacity;
  int count;
  Value* values;
} ValueArray;

void initValueArray(ValueArray* array);
void writeValueArray(ValueArray* array, Value value);
void freeValueArray(ValueArray* array);
void printValue(Value value);
bool valuesEqual(Value a, Value b);
#endif