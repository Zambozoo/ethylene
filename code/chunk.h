#ifndef clox_chunk_h
#define clox_chunk_h

#include "common.h"
#include "value.h"

#define STACK_MAX 256 //TODO: Grow stack dynamically

typedef enum {
  OP_CONSTANT_8, //Loads in constant
  OP_CONSTANT_16, //Loads in constant
  OP_NULL,
  OP_TRUE,
  OP_FALSE,

  //Undecided until type check
  OP_RETURN,
  OP_NEGATE,
  OP_ADD,
  OP_SUBTRACT,
  OP_MULTIPLY,
  OP_DIVIDE,
  OP_EQUAL,
  OP_GREATER,
  OP_LESS,

  //OTHER OPERATIONS
  OP_REFERENCE,
  OP_DEREFERENCE,

  OP_ASSIGN,
  OP_TYPE,

  OP_PRIMITIVE_HASH,
  OP_PRIMITIVE_EQ,
  OP_PRIMITIVE_BANG_EQ,

  #pragma region Integer Ops
  OP_INT_TO_CHAR,
  OP_INT_TO_BOOL,
  OP_INT_TO_FLT,
  OP_INT_8_BIT,
  OP_INT_16_BIT,
  OP_INT_32_BIT,
  OP_INT_64_BIT,

  OP_INT_NEGATE,
  OP_INT_INC_PREFIX,
  OP_INT_DEC_PREFIX,
  OP_INT_INC_POSTFIX,
  OP_INT_DEC_POSTFIX,

  OP_INT_ADD,
  OP_INT_SUBTRACT,
  OP_INT_MULTIPLY,
  OP_INT_DIVIDE,
  OP_INT_MODULO,

  OP_INT_LESSER,
  OP_INT_GREATER,
  OP_INT_LESSER_EQ,
  OP_INT_GREATER_EQ,
  OP_INT_COMPARE,

  OP_INT_RETURN,
  #pragma endregion

  #pragma region Float Ops
  OP_FLT_TO_INT,
  OP_FLT_8_TO_64,
  OP_FLT_16_TO_64,
  OP_FLT_32_O_64,
  OP_FLT_64_TO_8,
  OP_FLT_64_TO_16,
  OP_FLT_64_TO_32,

  OP_FLT_NEGATE,
  OP_FLT_INC_PREFIX,
  OP_FLT_DEC_PREFIX,
  OP_FLT_INC_POSTFIX,
  OP_FLT_DEC_POSTFIX,

  OP_FLT_ADD,
  OP_FLT_SUBTRACT,
  OP_FLT_MULTIPLY,
  OP_FLT_DIVIDE,

  OP_FLT_LESSER,
  OP_FLT_GREATER,
  OP_FLT_LESSER_EQ,
  OP_FLT_GREATER_EQ,
  OP_FLT_COMPARE, 
  
  OP_FLT_RETURN,
  #pragma endregion

  #pragma region Boolean Operations
  OP_BOOL_TO_INT,

  OP_BOOL_AND,
  OP_BOOL_OR,
  OP_BOOL_XOR,
  OP_BOOL_NEGATE,

  OP_BOOL_TERNARY,
  #pragma endregion

  #pragma region Character Operations
  OP_CHAR_TO_INT,
  OP_CHAR_INC_PREFIX,
  OP_CHAR_DEC_PREFIX,
  OP_CHAR_DEC_POSTFIX,

  OP_CHAR_INT_ADD,
  OP_CHAR_INT_SUBTRACT,
  OP_INT_CHAR_ADD,
  OP_INT_CHAR_SUBTRACT,

  OP_CHAR_LESSER,
  OP_CHAR_GREATER,
  OP_CHAR_LESSER_EQ,
  OP_CHAR_GREATER_EQ,
  OP_CHAR_COMPARE,
  #pragma endregion

  #pragma region Bit
  OP_BIT_AND,
  OP_BIT_XOR,
  OP_BIT_OR,
  OP_BIT_NEGATE,
  OP_BIT_RIGHT_SHIFT,
  OP_BIT_LEFT_SHIFT,
  #pragma endregion

} OpCode;

/**
 * Chunk
 * List of constants, 
 * list of op codes, 
 * associative array of line -> # of opCodes.
*/
typedef struct {
  int codeCount;
  int codeCapacity;
  uint8_t* code;

  int linesCapacity;
  uint16_t* lines;
  
  ValueArray constants;
} Chunk; //ArrayList

void initChunk(Chunk* chunk);
void freeChunk(Chunk* chunk);
void writeChunk(Chunk* chunk, uint8_t byte, int line);
int writeConstant(Chunk* chunk, Value value, int line);

int getLine(Chunk* chunk, int index);

#endif