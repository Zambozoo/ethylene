#include <stdio.h>
#include <stdint.h>
#include <stdarg.h>

#include "common.h"
#include "compiler.h"
#include "debug.h"
#include "vm.h"
#include "floats.h"
#include "value.h"

VM vm; 

//Empty stack
static void resetStack() {
  vm.stackTop = vm.stack;
}

static void runtimeError(const char* format, ...) {
  va_list args;
  va_start(args, format);
  vfprintf(stderr, format, args);
  va_end(args);
  fputs("\n", stderr);

  size_t instruction = vm.ip - vm.chunk->code - 1;
  int line = vm.chunk->lines[instruction];
  fprintf(stderr, "[line %d] in script\n", line);
  resetStack();
}

//Initialize the VM
void initVM() {
    resetStack();
}

//Free the VM
void freeVM() {
}

void push(Value value) {
  *vm.stackTop = value;
  vm.stackTop++;
}

Value pop() {
  vm.stackTop--;
  return *vm.stackTop;
}

static Value peek(int distance) {
  return vm.stackTop[-1 - distance];
}

void debugPrint() {
    printf("          ");
    for (Value* slot = vm.stack; slot < vm.stackTop; slot++) {
      printf("[ ");
      printValue(*slot);
      printf(" ]");
    }
    printf("\n");
    disassembleInstructionNoLine(vm.chunk, (int)(vm.ip - vm.chunk->code));
}

//Run the vm
static InterpretResult run() {
#define READ_BYTE() (*vm.ip++)
#define READ_SHORT() ((*vm.ip++) + (*vm.ip++ << 8))
#define READ_CONSTANT_8() (vm.chunk->constants.values[READ_BYTE()])
#define READ_CONSTANT_16() (vm.chunk->constants.values[READ_SHORT()])
#define BIN_UNTYPED_OP(isBool, op) do { \
  if (peek(0).type == peek(1).type) { \
    if(IS_FLT(peek(0))) { \
      double b = AS_FLT(pop()); \
      double a = AS_FLT(pop()); \
      if(isBool) push(BOOL_VAL(a op b)); \
      else push(FLT_VAL(a op b)); \
      break; \
    } else if(IS_INT(peek(0))) { \
      int64_t b = AS_INT(pop()); \
      int64_t a = AS_INT(pop()); \
      if(isBool) push(BOOL_VAL(a op b)); \
      else push(INT_VAL(a op b)); \
      break; \
    } \
  } \
  runtimeError("Operands must be ints/flts."); \
  return INTERPRET_RUNTIME_ERROR; \
} while (false)
#define BINARY_OP(typeCheck, typeCast, valueType, type, op) do { \
  if (!typeCheck(peek(0)) || !typeCheck(peek(1))) { \
    runtimeError("Operands must be ints/flts."); \
    return INTERPRET_RUNTIME_ERROR; \
  } \
  type b = typeCast(pop()); \
  type a = typeCast(pop()); \
  push(valueType(a op b)); \
} while (false)
#define SHIFT_OP(op) do { \
  if (!IS_BIT(peek(0)) || !IS_INT(peek(1))) { \
    runtimeError("Operands must be bit op int."); \
    return INTERPRET_RUNTIME_ERROR; \
  } \
  int64_t b = AS_INT(pop()); \
  bit64 a = AS_BIT(pop()); \
  push(BIT_VAL(a op b)); \
} while (false)
#define COMPARE(typeCheck, typeCast, valueType, type) do { \
  if (!typeCheck(peek(0)) || !typeCheck(peek(1))) { \
    runtimeError("Operands must be ints."); \
    return INTERPRET_RUNTIME_ERROR; \
  } \
  type b = typeCast(pop()); \
  type a = typeCast(pop()); \
  push(valueType((a > b) - (a < b))); \
} while (false)

  for (;;) {
#ifdef DEBUG_TRACE_EXECUTION
    debugPrint();
#endif
    uint8_t instruction;
    switch (instruction = READ_BYTE()) {
      case OP_CONSTANT_8: push(READ_CONSTANT_8()); break;
      case OP_CONSTANT_16: push(READ_CONSTANT_16()); break;
      case OP_EQUAL: {
        if(peek(0).type != peek(1).type) {
          runtimeError("Operands must be the same type.");
          return INTERPRET_RUNTIME_ERROR;
        }
        Value b = pop();
        Value a = pop();
        push(BOOL_VAL(valuesEqual(a, b)));
        break;
      }
      case OP_ADD:      BIN_UNTYPED_OP(false, +); break;
      case OP_SUBTRACT: BIN_UNTYPED_OP(false, -); break;
      case OP_MULTIPLY: BIN_UNTYPED_OP(false, *); break;
      case OP_DIVIDE:   BIN_UNTYPED_OP(false, /); break;
      case OP_GREATER:  BIN_UNTYPED_OP(true, >); break;
      case OP_LESS:     BIN_UNTYPED_OP(true, <); break;
    #pragma region Integer
    //   case OP_INT_TO_CHAR:
      case OP_INT_TO_BOOL: push(BOOL_VAL(!!AS_INT(pop()))); break;
      case OP_INT_TO_FLT: { double d = AS_INT(pop()); push(FLT_VAL(d)); break; }
      case OP_INT_8_BIT:  push(INT_VAL(0xff & AS_INT(pop()))); break;
      case OP_INT_16_BIT: push(INT_VAL(0xffff & AS_INT(pop()))); break;
      case OP_INT_32_BIT: push(INT_VAL(0xffffffff & AS_INT(pop()))); break;

      case OP_INT_NEGATE: push(INT_VAL(-AS_INT(pop()))); break;
    //   case OP_INT_INC_PREFIX:
    //   case OP_INT_DEC_PREFIX:
    //   case OP_INT_INC_POSTFIX:
    //   case OP_INT_DEC_POSTFIX:

      case OP_INT_ADD:          BINARY_OP(IS_INT, AS_INT, INT_VAL, int64_t, +); break;
      case OP_INT_SUBTRACT:     BINARY_OP(IS_INT, AS_INT, INT_VAL, int64_t, -); break;
      case OP_INT_MULTIPLY:     BINARY_OP(IS_INT, AS_INT, INT_VAL, int64_t, *); break;
      case OP_INT_DIVIDE:       BINARY_OP(IS_INT, AS_INT, INT_VAL, int64_t, /); break;
      case OP_INT_MODULO:       BINARY_OP(IS_INT, AS_INT, INT_VAL, int64_t, %); break;

      case OP_INT_LESSER:       BINARY_OP(IS_INT, AS_INT, BOOL_VAL, int64_t, <); break;
      case OP_INT_GREATER:      BINARY_OP(IS_INT, AS_INT, BOOL_VAL, int64_t, >); break;
      case OP_INT_LESSER_EQ:    BINARY_OP(IS_INT, AS_INT, BOOL_VAL, int64_t, <=); break;
      case OP_INT_GREATER_EQ:   BINARY_OP(IS_INT, AS_INT, BOOL_VAL, int64_t, >=); break;
      case OP_INT_COMPARE:      COMPARE(IS_INT, AS_INT, INT_VAL, int64_t); break;
    #pragma endregion

    #pragma region Float
      case OP_FLT_TO_INT:   { int64_t i = AS_FLT(pop()); push(INT_VAL(i)); break; }
      case OP_FLT_8_TO_64:  push(FLT_VAL(quarterToFloat((float8){.u=AS_INT(pop())}).f)); break;
      case OP_FLT_16_TO_64: push(FLT_VAL(halfToFloat((float16){.u=AS_INT(pop())}).f)); break;
      case OP_FLT_32_O_64:  push(FLT_VAL((float32){.u=AS_INT(pop())}.f)); break;
      case OP_FLT_64_TO_8:  { push(INT_VAL(floatToQuarter((float32){.f=AS_FLT(pop())}).u)); break; }
      case OP_FLT_64_TO_16: { push(INT_VAL(floatToHalf((float32){.f=AS_FLT(pop())}).u)); break; }
      case OP_FLT_64_TO_32: { push(INT_VAL((float32){.f=AS_FLT(pop())}.u)); break; }

      case OP_FLT_NEGATE: 
        if (!IS_FLT(peek(0))) {
          runtimeError("Operand must be a number.");
          return INTERPRET_RUNTIME_ERROR;
        }
        push(FLT_VAL(-AS_FLT(pop()))); break;
    //   case OP_FLT_INC_PREFIX:
    //   case OP_FLT_DEC_PREFIX:
    //   case OP_FLT_INC_POSTFIX:
    //   case OP_FLT_DEC_POSTFIX:

      case OP_FLT_ADD:      BINARY_OP(IS_FLT, AS_FLT, FLT_VAL, double,+); break;
      case OP_FLT_SUBTRACT: BINARY_OP(IS_FLT, AS_FLT, FLT_VAL, double, -); break;
      case OP_FLT_MULTIPLY: BINARY_OP(IS_FLT, AS_FLT, FLT_VAL, double, *); break;
      case OP_FLT_DIVIDE:   BINARY_OP(IS_FLT, AS_FLT, FLT_VAL, double, /); break;

      case OP_FLT_LESSER:       BINARY_OP(IS_FLT, AS_FLT, BOOL_VAL, double, <); break;
      case OP_FLT_GREATER:      BINARY_OP(IS_FLT, AS_FLT, BOOL_VAL, double, >); break;
      case OP_FLT_LESSER_EQ:    BINARY_OP(IS_FLT, AS_FLT, BOOL_VAL, double, <=); break;
      case OP_FLT_GREATER_EQ:   BINARY_OP(IS_FLT, AS_FLT, BOOL_VAL, double, >=); break;
      case OP_FLT_COMPARE:      COMPARE(IS_FLT, AS_FLT, FLT_VAL, double); break;
    #pragma endregion

    
    #pragma region Bit
    case OP_BIT_AND:          BINARY_OP(IS_BIT, AS_BIT, BIT_VAL, bit64, &); break;
    case OP_BIT_OR:           BINARY_OP(IS_BIT, AS_BIT, BIT_VAL, bit64, |); break;
    case OP_BIT_XOR:          BINARY_OP(IS_BIT, AS_BIT, BIT_VAL, bit64, ^); break;
    case OP_BIT_RIGHT_SHIFT:  SHIFT_OP(>>); break;
    case OP_BIT_LEFT_SHIFT:   SHIFT_OP(<<); break;
    case OP_BIT_NEGATE: push(INT_VAL(~AS_INT(pop()))); break;
    #pragma endregion
      case OP_NULL: push(NULL_VAL); break;
      case OP_TRUE: push(BOOL_VAL(true)); break;
      case OP_FALSE: push(BOOL_VAL(false)); break;

      
      case OP_BOOL_NEGATE:
        if (!IS_BOOL(peek(0))) {
          runtimeError("Operand must be bool.");
          return INTERPRET_RUNTIME_ERROR;
        }
        push(BOOL_VAL(!AS_BOOL(pop())));
        break;

      case OP_INT_RETURN: {
        printValue(pop());
        printf("\n");
        return INTERPRET_OK;
      }
      case OP_FLT_RETURN: {
        printValue(pop());
        printf("\n");
        return INTERPRET_OK;
      }
      case OP_RETURN: {
        printValue(pop());
        printf("\n");
        return INTERPRET_OK;
      }
    }
  }

#undef READ_BYTE
#undef READ_SHORT
#undef READ_CONSTANT_8
#undef READ_CONSTANT_16
#undef BINARY_OP
#undef CYCLE_N
#undef COMPARE
}

InterpretResult interpret(const char* source) {
  Chunk chunk;
  initChunk(&chunk);

  if (!compile(source, &chunk)) {
    freeChunk(&chunk);
    return INTERPRET_COMPILE_ERROR;
  }

  vm.chunk = &chunk;
  vm.ip = vm.chunk->code;

  InterpretResult result = run();

  freeChunk(&chunk);
  return result;
}