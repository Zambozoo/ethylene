#ifndef clox_debug_h
#define clox_debug_h

#include "chunk.h"

void disassembleChunk(Chunk* chunk, const char* name);
int disassembleInstructionNoLine(Chunk* chunk, int offset);
int disassembleInstruction(Chunk* chunk, int offset, int line, int oldLine);

#endif