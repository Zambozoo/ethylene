#include <stdio.h>

#include "debug.h"

void disassembleChunk(Chunk* chunk, const char* name) {
  printf("== %s ==\n", name);

  for (int line = 0, offset = 0; offset < chunk->codeCount;) {
    int oldLine = line;
    line = getLine(chunk, offset);
    offset = disassembleInstruction(chunk, offset, line, oldLine);
  }
}

static int constant8Instruction(const char* name, Chunk* chunk, int offset) {
  uint8_t constant = chunk->code[offset + 1];
  printf("%-16s %4d '", name, constant);
  printValue(chunk->constants.values[constant]);
  printf("'\n");
  return offset + 2;
}

static int constant16Instruction(const char* name, Chunk* chunk, int offset) {
  uint16_t constant = (chunk->code[offset + 1] << 8) | chunk->code[offset + 2];
  printf("%-16s %4d '", name, constant);
  printValue(chunk->constants.values[constant]);
  printf("'\n");
  return offset + 3;
}

static int simpleInstruction(const char* name, int offset) {
  printf("%s\n", name);
  return offset + 1;
}

int disassembleInstructionNoLine(Chunk* chunk, int offset) {
  return disassembleInstruction(chunk, offset, 0, 0);
} 

int disassembleInstruction(Chunk* chunk, int offset, int line, int oldLine) {
  printf("%04d ", offset);
  
  if (line == oldLine) {
    printf("   | ");
  } else {
    printf("%4d ", line);
  }

  uint8_t instruction = chunk->code[offset];
  switch (instruction) {
    case OP_CONSTANT_8: return constant8Instruction("OP_CONSTANT_8", chunk, offset);
    case OP_CONSTANT_16: return constant16Instruction("OP_CONSTANT_16", chunk, offset);

  #pragma Untyped
    case OP_RETURN: return simpleInstruction("OP_RETURN", offset);
    case OP_NEGATE: return simpleInstruction("OP_NEGATE", offset);
    case OP_ADD: return simpleInstruction("OP_ADD", offset);
    case OP_SUBTRACT: return simpleInstruction("OP_SUBTRACT", offset);
    case OP_MULTIPLY: return simpleInstruction("OP_MULTIPLY", offset);
    case OP_DIVIDE: return simpleInstruction("OP_DIVIDE", offset);
    case OP_EQUAL: return simpleInstruction("OP_EQUAL", offset);
    case OP_GREATER: return simpleInstruction("OP_GREATER", offset);
    case OP_LESS: return simpleInstruction("OP_LESS", offset);
  #pragma endregion

  #pragma region Integer
    // case OP_INT_TO_CHAR: return simpleInstruction("OP_INT_TO_CHAR", offset);
    // case OP_INT_TO_BOOL: return simpleInstruction("OP_INT_TO_BOOL", offset);
    // case OP_INT_TO_FLT: return simpleInstruction("OP_INT_TO_FLT", offset);
    case OP_INT_8_BIT: return simpleInstruction("OP_INT_8_BIT", offset);
    case OP_INT_16_BIT: return simpleInstruction("OP_INT_16_BIT", offset);
    case OP_INT_32_BIT: return simpleInstruction("OP_INT_32_BIT", offset);
    case OP_INT_64_BIT: return simpleInstruction("OP_INT_64_BIT", offset);

    case OP_INT_NEGATE: return simpleInstruction("OP_INT_NEGATE", offset);
    // case OP_INT_INC_PREFIX: return simpleInstruction("OP_INT_INC_PREFIX", offset);
    // case OP_INT_DEC_PREFIX: return simpleInstruction("OP_INT_DEC_PREFIX", offset);
    // case OP_INT_INC_POSTFIX: return simpleInstruction("OP_INT_INC_POSTFIX", offset);
    // case OP_INT_DEC_POSTFIX: return simpleInstruction("OP_INT_DEC_POSTFIX", offset);

    case OP_INT_ADD: return simpleInstruction("OP_INT_ADD", offset);
    case OP_INT_SUBTRACT: return simpleInstruction("OP_INT_SUBTRACT", offset);
    case OP_INT_MULTIPLY: return simpleInstruction("OP_INT_MULTIPLY", offset);
    case OP_INT_DIVIDE: return simpleInstruction("OP_INT_DIVIDE", offset);
    case OP_INT_MODULO: return simpleInstruction("OP_INT_MODULO", offset);

    case OP_INT_LESSER: return simpleInstruction("OP_INT_LESSER", offset);
    case OP_INT_GREATER: return simpleInstruction("OP_INT_GREATER", offset);
    case OP_INT_LESSER_EQ: return simpleInstruction("OP_INT_LESSER_EQ", offset);
    case OP_INT_GREATER_EQ: return simpleInstruction("OP_INT_GREATER_EQ", offset);
    case OP_INT_COMPARE: return simpleInstruction("OP_INT_COMPARE", offset);
    
    case OP_INT_RETURN: return simpleInstruction("OP_INT_RETURN", offset);
  #pragma endregion
    
  #pragma region Float
    case OP_FLT_TO_INT: return simpleInstruction("OP_FLT_TO_INT", offset);
    case OP_FLT_8_TO_64: return simpleInstruction("OP_FLT_8_TO_64", offset);
    case OP_FLT_16_TO_64: return simpleInstruction("OP_FLT_16_TO_64", offset);
    case OP_FLT_32_O_64: return simpleInstruction("OP_FLT_32_O_64", offset);
    case OP_FLT_64_TO_8: return simpleInstruction("OP_FLT_64_TO_8", offset);
    case OP_FLT_64_TO_16: return simpleInstruction("OP_FLT_64_TO_16", offset);
    case OP_FLT_64_TO_32: return simpleInstruction("OP_FLT_64_TO_32", offset);

    case OP_FLT_NEGATE: return simpleInstruction("OP_FLT_NEGATE", offset);
    case OP_FLT_INC_PREFIX: return simpleInstruction("OP_FLT_INC_PREFIX", offset);
    case OP_FLT_DEC_PREFIX: return simpleInstruction("OP_FLT_DEC_PREFIX", offset);
    case OP_FLT_INC_POSTFIX: return simpleInstruction("OP_FLT_INC_POSTFIX", offset);
    case OP_FLT_DEC_POSTFIX: return simpleInstruction("OP_FLT_DEC_POSTFIX", offset);

    case OP_FLT_ADD: return simpleInstruction("OP_FLT_ADD", offset);
    case OP_FLT_SUBTRACT: return simpleInstruction("OP_FLT_SUBTRACT", offset);
    case OP_FLT_MULTIPLY: return simpleInstruction("OP_FLT_MULTIPLY", offset);
    case OP_FLT_DIVIDE: return simpleInstruction("OP_FLT_DIVIDE", offset);

    case OP_FLT_LESSER: return simpleInstruction("OP_FLT_LESSER", offset);
    case OP_FLT_GREATER: return simpleInstruction("OP_FLT_GREATER", offset);
    case OP_FLT_LESSER_EQ: return simpleInstruction("OP_FLT_LESSER_EQ", offset);
    case OP_FLT_GREATER_EQ: return simpleInstruction("OP_FLT_GREATER_EQ", offset);
    case OP_FLT_COMPARE: return simpleInstruction("OP_FLT_COMPARE", offset);

    case OP_FLT_RETURN: return simpleInstruction("OP_FLT_RETURN", offset);
  #pragma endregion

  #pragma region Boolean
    case OP_BOOL_TO_INT: return simpleInstruction("OP_BOOL_TO_INT", offset);
    case OP_BOOL_AND: return simpleInstruction("OP_BOOL_AND", offset);
    case OP_BOOL_OR: return simpleInstruction("OP_BOOL_OR", offset);
    case OP_BOOL_XOR: return simpleInstruction("OP_BOOL_XOR", offset);
    case OP_BOOL_TERNARY: return simpleInstruction("OP_BOOL_TERNARY", offset);
    case OP_BOOL_NEGATE:  simpleInstruction("OP_BOOL_NEGATE", offset);
  #pragma endregion

  #pragma region Bit
    case OP_BIT_AND: return simpleInstruction("OP_BIT_AND", offset);
    case OP_BIT_OR: return simpleInstruction("OP_BIT_OR", offset);
    case OP_BIT_XOR: return simpleInstruction("OP_BIT_XOR", offset);
    case OP_BIT_NEGATE: return simpleInstruction("OP_BIT_NEGATE", offset);
    
    case OP_BIT_RIGHT_SHIFT: return simpleInstruction("OP_BIT_RIGHT_SHIFT", offset);
    case OP_BIT_LEFT_SHIFT: return simpleInstruction("OP_BIT_LEFT_SHIFT", offset);
  #pragma endregion

    case OP_NULL: return simpleInstruction("OP_NULL", offset);
    case OP_TRUE: return simpleInstruction("OP_TRUE", offset);
    case OP_FALSE: return simpleInstruction("OP_FALSE", offset);
    default:
      printf("Unknown opcode %d\n", instruction);
      return offset + 1;
  }
}