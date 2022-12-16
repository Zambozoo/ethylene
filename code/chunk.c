#include <stdlib.h>
#include <stdio.h>
#include "chunk.h"
#include "memory.h"

void initChunk(Chunk* chunk) {
  chunk->codeCount = 0;
  chunk->codeCapacity = 0;
  chunk->code = NULL;
  
  chunk->linesCapacity = 0;
  chunk->lines = NULL;
  initValueArray(&chunk->constants);
}

void freeChunk(Chunk* chunk) {
  FREE_ARRAY(uint8_t, chunk->code, chunk->codeCapacity);
  FREE_ARRAY(int, chunk->lines, chunk->linesCapacity);
  freeValueArray(&chunk->constants);
  initChunk(chunk);
}

void writeChunk(Chunk* chunk, uint8_t byte, int line) {
  if (chunk->codeCapacity < chunk->codeCount + 1) {
    int oldCapacity = chunk->codeCapacity;
    chunk->codeCapacity = GROW_CAPACITY(oldCapacity);
    chunk->code = GROW_ARRAY(uint8_t, chunk->code, oldCapacity, chunk->codeCapacity);
  }

  if (chunk->linesCapacity < line) {
    int oldCapacity;
    do {
      oldCapacity = chunk->linesCapacity;
      chunk->linesCapacity = GROW_CAPACITY(oldCapacity);
    } while(chunk->linesCapacity < line);
    chunk->lines = GROW_ARRAY(uint16_t, chunk->lines, oldCapacity, chunk->linesCapacity);
    for(; oldCapacity < chunk->linesCapacity; oldCapacity++) {
      chunk->lines[oldCapacity] = 0;
    }
  }
  chunk->lines[line - 1]++;
  chunk->code[chunk->codeCount] = byte;
  chunk->codeCount++;
}

int writeConstant(Chunk* chunk, Value i, int line) {
  writeValueArray(&chunk->constants, i);
  if(chunk->constants.count < 128) {
    writeChunk(chunk, OP_CONSTANT_8, line);
  } else {
    writeChunk(chunk, OP_CONSTANT_16, line);
  }
  return chunk->constants.count - 1;
}

int getLine(Chunk* chunk, int index) {
  uint16_t* line = chunk->lines;
  int lineLength = chunk->lines[0];
  for (int offset = 0; offset <= index; offset++) {
    while(lineLength == 0) {
      lineLength = *(++line);
    }
    lineLength--;
  }
  return (line - chunk->lines) + 1;
}