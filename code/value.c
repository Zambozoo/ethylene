#include <stdio.h>

#include "memory.h"
#include "value.h"

void initValueArray(ValueArray* array) {
  array->values = NULL;
  array->capacity = 0;
  array->count = 0;
}

void writeValueArray(ValueArray* array, Value value) {
  if (array->capacity < array->count + 1) {
    int oldCapacity = array->capacity;
    array->capacity = GROW_CAPACITY(oldCapacity);
    array->values = GROW_ARRAY(Value, array->values, oldCapacity, array->capacity);
  }
  array->values[array->count] = value;
  array->count++;
}

void freeValueArray(ValueArray* array) {
  FREE_ARRAY(Value, array->values, array->capacity);
  initValueArray(array);
}

static uint64_t bits(uint64_t n)
{
    uint64_t m = n ? bits(n / 2) : 0;
    printf("%d", (int)(n % 2));
    return m;
}

void printValue(Value value) {
  switch (value.type) {
    case VAL_BOOL: printf(AS_BOOL(value) ? "true" : "false"); break;
    case VAL_NULL: printf("null"); break;
    case VAL_FLT: printf("%g", AS_FLT(value)); break;
    case VAL_INT: printf("%ld ", AS_INT(value)); break;
    case VAL_BIT: bits(AS_BIT(value)); break;
  }
}

bool valuesEqual(Value a, Value b) {
  if (a.type != b.type) return false;
  switch (a.type) {
    case VAL_BOOL:  return AS_BOOL(a) == AS_BOOL(b);
    case VAL_NULL:  return true;
    case VAL_INT:   return AS_INT(a) == AS_INT(b);
    case VAL_FLT:   return AS_FLT(a) == AS_FLT(b);
    case VAL_BIT:   return AS_BIT(a) == AS_BIT(b);
    default:        return false; // Unreachable.
  }
}