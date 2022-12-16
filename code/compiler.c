#include <stdio.h>
#include <stdlib.h>

#include "common.h"
#include "compiler.h"
#include "scanner.h"
#include "value.h"

#ifdef DEBUG_PRINT_CODE
#include "debug.h"
#endif

typedef struct {
  Token current;
  Token previous;
  bool hadError;
  bool panicMode;
} Parser;

typedef enum {
  PREC_NONE,
  PREC_ASSIGNMENT,  // = ?: ??
  PREC_OR,          // ||
  PREC_AND,         // &&
  PREC_BIT_OR,      // |
  PREC_BIT_XOR,     // ^
  PREC_BIT_AND,     // &
  PREC_EQUALITY,    // == !=
  PREC_COMPARISON,  // < > <= >= 
  PREC_THREE_CMP,   // <=>
  PREC_SHIFT,       // << >>
  PREC_TERM,        // + -
  PREC_FACTOR,      // * / %
  PREC_UNARY_POST,  // ++ -- * & 
  PREC_UNARY_PRE,   // ! - ++ -- @ #
  PREC_CALL,        // . () [] {}
  PREC_PRIMARY
} Precedence;

typedef void (*ParseFn)();

//TODO: ADD POSTFIX
//TODO: ADD TERNARY
//TODO: FIGURE OUT TYPES BRUH
typedef struct {
  ParseFn prefix;
  ParseFn infix;
  Precedence precedence;
} ParseRule;

Parser parser;

Chunk* compilingChunk;

static Chunk* currentChunk() {
  return compilingChunk;
}

#pragma region Parser

static void errorAt(Token* token, const char* message) {
  if (parser.panicMode) return;
  parser.panicMode = true;
  fprintf(stderr, "[line %d] Error", token->line);

  if (token->type == TOKEN_EOF) {
    fprintf(stderr, " at end");
  } else if (token->type == TOKEN_ERROR) {
    // Nothing.
  } else {
    fprintf(stderr, " at '%.*s'", token->length, token->start);
  }

  fprintf(stderr, ": %s\n", message);
  parser.hadError = true;
}

static void error(const char* message) {
  errorAt(&parser.previous, message);
}

static void errorAtCurrent(const char* message) {
  errorAt(&parser.current, message);
}

static void advance() {
  parser.previous = parser.current;

  for (;;) {
    parser.current = scanToken();
    if (parser.current.type != TOKEN_ERROR) break;

    errorAtCurrent(parser.current.start);
  }
}

static void consume(TokenType type, const char* message) {
  if (parser.current.type == type) {
    advance();
    return;
  }

  errorAtCurrent(message);
}

static void emitByte(uint8_t byte) {
  writeChunk(currentChunk(), byte, parser.previous.line);
}

static void emitBytes(uint8_t byte1, uint8_t byte2) {
  printf("BYTE1:%d,BYTE2:%d\n",byte1, byte2);
  emitByte(byte1);
  emitByte(byte2);
}

static void emitReturn() {
  emitByte(OP_RETURN);
}

static uint16_t makeConstant(Value value, int line) {
  int constant = writeConstant(currentChunk(), value, line);
  if (constant > UINT16_MAX) {
    error("Too many constants in one chunk.");
    return 0;
  }

  return (uint16_t)constant;
}

static void emitConstant(Value value, int line) {
  uint16_t c = makeConstant(value, line);
  emitByte(c);
}

static void endCompiler() {
  emitReturn();

#ifdef DEBUG_PRINT_CODE
  if (!parser.hadError) {
    disassembleChunk(currentChunk(), "code");
  }
#endif
}

#pragma endregion

#pragma region Expressions
static void expression();
static ParseRule* getRule(TokenType type);
static void parsePrecedence(Precedence precedence);

static void binary() {
  TokenType operatorType = parser.previous.type;
  ParseRule* rule = getRule(operatorType);
  parsePrecedence((Precedence)(rule->precedence + 1));

  switch (operatorType) {
    case TOKEN_PLUS:          emitByte(OP_ADD); break;
    case TOKEN_MINUS:         emitByte(OP_SUBTRACT); break;
    case TOKEN_STAR:          emitByte(OP_MULTIPLY); break;
    case TOKEN_SLASH:         emitByte(OP_DIVIDE); break;
    case TOKEN_MODULO:        emitByte(OP_INT_MODULO); break;

    case TOKEN_BIT_XOR:       emitByte(OP_BIT_XOR); break;
    case TOKEN_BIT_AND:       emitByte(OP_BIT_AND); break;
    case TOKEN_BIT_OR:        emitByte(OP_BIT_OR); break;

    case TOKEN_R_SHIFT:       emitByte(OP_BIT_RIGHT_SHIFT); break;
    case TOKEN_L_SHIFT:       emitByte(OP_BIT_LEFT_SHIFT); break;
    
    case TOKEN_OR:            emitByte(OP_BOOL_OR); break;
    case TOKEN_AND:           emitByte(OP_BOOL_AND); break;

    case TOKEN_BANG_EQUAL:    emitBytes(OP_EQUAL, OP_BOOL_NEGATE); break;
    case TOKEN_EQUAL_EQUAL:   emitByte(OP_EQUAL); break;

    case TOKEN_GREATER:       emitByte(OP_GREATER); break;
    case TOKEN_GREATER_EQUAL: emitBytes(OP_LESS, OP_BOOL_NEGATE); break;
    case TOKEN_LESS:          emitByte(OP_LESS); break;
    case TOKEN_LESS_EQUAL:    emitBytes(OP_GREATER, OP_BOOL_NEGATE); break;
    default: return; // Unreachable.
  }
}

static void literal() {
  switch (parser.previous.type) {
    case TOKEN_FALSE: emitByte(OP_FALSE); break;
    case TOKEN_NULL: emitByte(OP_NULL); break;
    case TOKEN_TRUE: emitByte(OP_TRUE); break;
    default: return; // Unreachable.
  }
}

static void grouping() {
  expression();
  consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.");
}

static void _float(int line) {
  double d = strtod(parser.previous.start, NULL);
  emitConstant(FLT_VAL(d), line);
}

static void _int(int line) {
  int64_t i = strtoll(parser.previous.start, NULL, 10);
  emitConstant(INT_VAL(i), line);
}

static void unary() {
  TokenType operatorType = parser.previous.type;

  // Compile the operand.
  parsePrecedence(PREC_UNARY_PRE);

  // Emit the operator instruction.
  switch (operatorType) {
    case TOKEN_BANG: emitByte(OP_BOOL_NEGATE); break;
    case TOKEN_MINUS: emitByte(OP_NEGATE); break;
    default: return; // Unreachable.
  }
}

static void ternary() {
  parsePrecedence(PREC_ASSIGNMENT);
  consume(TOKEN_COLON, "expected ':' in ternary");
  parsePrecedence(PREC_ASSIGNMENT);
}

ParseRule rules[] = {
  [TOKEN_QMARK]         = {NULL,     ternary,PREC_ASSIGNMENT},
  [TOKEN_COLON]         = {NULL,     NULL,   PREC_NONE},
  [TOKEN_LEFT_PAREN]    = {grouping, NULL,   PREC_NONE},
  [TOKEN_RIGHT_PAREN]   = {NULL,     NULL,   PREC_NONE},
  [TOKEN_LEFT_BRACE]    = {NULL,     NULL,   PREC_NONE}, 
  [TOKEN_RIGHT_BRACE]   = {NULL,     NULL,   PREC_NONE},
  [TOKEN_COMMA]         = {NULL,     NULL,   PREC_NONE},
  [TOKEN_DOT]           = {NULL,     NULL,   PREC_NONE},
  [TOKEN_MINUS]         = {unary,    binary, PREC_TERM},
  [TOKEN_PLUS]          = {NULL,     binary, PREC_TERM},
  [TOKEN_SEMICOLON]     = {NULL,     NULL,   PREC_NONE},
  [TOKEN_SLASH]         = {NULL,     binary, PREC_FACTOR},
  [TOKEN_STAR]          = {NULL,     binary, PREC_FACTOR},
  [TOKEN_BANG]          = {unary,    NULL,   PREC_NONE},
  [TOKEN_BANG_EQUAL]    = {NULL,     binary, PREC_EQUALITY},
  [TOKEN_EQUAL]         = {NULL,     NULL,   PREC_NONE},
  [TOKEN_EQUAL_EQUAL]   = {NULL,     binary, PREC_EQUALITY},
  [TOKEN_GREATER]       = {NULL,     binary, PREC_COMPARISON},
  [TOKEN_GREATER_EQUAL] = {NULL,     binary, PREC_COMPARISON},
  [TOKEN_LESS]          = {NULL,     binary, PREC_COMPARISON},
  [TOKEN_LESS_EQUAL]    = {NULL,     binary, PREC_COMPARISON},
  [TOKEN_IDENTIFIER]    = {NULL,     NULL,   PREC_NONE},

  [TOKEN_STRING]        = {NULL,     NULL,   PREC_NONE},
  [TOKEN_INT]           = {_int,     NULL,   PREC_NONE},
  [TOKEN_FLT]           = {_float,   NULL,   PREC_NONE},
  [TOKEN_FIX]           = {NULL,     NULL,   PREC_NONE},
  [TOKEN_CHAR]          = {NULL,     NULL,   PREC_NONE},

  [TOKEN_AND]           = {NULL,     NULL,   PREC_NONE},
  [TOKEN_CLASS]         = {NULL,     NULL,   PREC_NONE},
  [TOKEN_ELSE]          = {NULL,     NULL,   PREC_NONE},
  [TOKEN_FALSE]         = {literal,  NULL,   PREC_NONE},
  [TOKEN_FOR]           = {NULL,     NULL,   PREC_NONE},
  [TOKEN_FUN]           = {NULL,     NULL,   PREC_NONE},
  [TOKEN_IF]            = {NULL,     NULL,   PREC_NONE},
  [TOKEN_NULL]          = {literal,  NULL,   PREC_NONE},
  [TOKEN_OR]            = {NULL,     NULL,   PREC_NONE},
  [TOKEN_PRINT]         = {NULL,     NULL,   PREC_NONE},
  [TOKEN_RETURN]        = {NULL,     NULL,   PREC_NONE},
  [TOKEN_SUPER]         = {NULL,     NULL,   PREC_NONE},
  [TOKEN_THIS]          = {NULL,     NULL,   PREC_NONE},
  [TOKEN_TRUE]          = {literal,  NULL,   PREC_NONE},
  [TOKEN_VAR]           = {NULL,     NULL,   PREC_NONE},
  [TOKEN_WHILE]         = {NULL,     NULL,   PREC_NONE},
  [TOKEN_ERROR]         = {NULL,     NULL,   PREC_NONE},
  [TOKEN_EOF]           = {NULL,     NULL,   PREC_NONE},
};

static void parsePrecedence(Precedence precedence) {
  advance();
  ParseFn prefixRule = getRule(parser.previous.type)->prefix;
  if (prefixRule == NULL) {
    error("Expect expression.");
    return;
  }

  prefixRule();
  while (precedence <= getRule(parser.current.type)->precedence) {
    advance();
    ParseFn infixRule = getRule(parser.previous.type)->infix;
    infixRule();
  }
}


static ParseRule* getRule(TokenType type) {
  return &rules[type];
}

static void expression() {
  parsePrecedence(PREC_ASSIGNMENT);
}

#pragma endregion

bool compile(const char* source, Chunk* chunk) {
  initScanner(source);
  compilingChunk = chunk;
  parser.hadError = false;
  parser.panicMode = false;

  advance();
  expression();

  consume(TOKEN_EOF, "Expect end of expression.");
  endCompiler();
  return !parser.hadError;

  // int line = -1;
  // for (;;) {
  //   Token token = scanToken();
  //   if (token.line != line) {
  //     printf("%4d ", token.line);
  //     line = token.line;
  //   } else {
  //     printf("   | ");
  //   }
  //   printf("%2d '%.*s'\n", token.type, token.length, token.start); 

  //   if (token.type == TOKEN_EOF) break;
  // }
}