#ifndef clox_vm_h
#define clox_vm_h

#include "chunk.h"
#include "floats.h"

typedef struct {
  Chunk* chunk;
  uint8_t* ip;
  
  Value stack[STACK_MAX];
  Value* stackTop;
} VM; //Executes chunks of code

typedef enum {
  INTERPRET_OK,
  INTERPRET_COMPILE_ERROR,
  INTERPRET_RUNTIME_ERROR
} InterpretResult;

void initVM();
void freeVM();
InterpretResult interpret(const char* source);

//Stack methods
void push(Value value);
Value pop();

#endif