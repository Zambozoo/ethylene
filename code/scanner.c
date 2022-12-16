#include <stdio.h>
#include <string.h>

#include "common.h"
#include "scanner.h"

typedef struct {
  const char* start;
  const char* current;
  int line;
} Scanner;

Scanner scanner;

void initScanner(const char* source) {
  scanner.start = source;
  scanner.current = source;
  scanner.line = 1;
}

static bool isAlpha(char c) {
  return (c >= 'a' && c <= 'z') ||
         (c >= 'A' && c <= 'Z') ||
          c == '_';
}

static bool isDigit(char c) {
  return c >= '0' && c <= '9';
}

static bool isHex(char c) {
  return c >= 'A' && c <= 'F' || isDigit(c);
}

#pragma region Scanner

static bool isAtEnd() {
  return *scanner.current == '\0';
}

static char advance() {
  scanner.current++;
  return scanner.current[-1];
}

static char advanceN(int n) {
  scanner.current += n;
  return scanner.current[-1];
}

static char peek() {
  return *scanner.current;
}

static char peekNext() {
  if (isAtEnd()) return '\0';
  return scanner.current[1];
}

static char peekN(int n) {
  for(int i = 0; i < n; i++) if(scanner.current[i] == '\0') return '\0';
  return scanner.current[n];
}

static bool match(char expected) {
  if (isAtEnd()) return false;
  if (*scanner.current != expected) return false;
  scanner.current++;
  return true;
}

//Returns Token from scanner.start to scanner.current of TokenType type
static Token makeToken(TokenType type) {
  Token token;
  token.type = type;
  token.start = scanner.start;
  token.length = (int)(scanner.current - scanner.start);
  token.line = scanner.line;
  return token;
}

static Token errorToken(const char* message) {
  Token token;
  token.type = TOKEN_ERROR;
  token.start = message;
  token.length = (int)strlen(message);
  token.line = scanner.line;
  return token;
}

#pragma endregion

static void skipWhitespace() {
  for (;;) {
    char c = peek();
    switch (c) {
      case ' ': case '\r': case '\t': advance(); break;
      case '\n': scanner.line++; advance(); break;
      case '/':
        if (peekNext() == '/') {
          // A comment goes until the end of the line.
          while (peek() != '\n' && !isAtEnd()) advance();
        } else {
          return;
        }
        break;
      default: return;
    }
  }
}

#pragma region Keywords/Identifiers

static TokenType checkKeyword(int start, int length,
    const char* rest, TokenType type) {
  if (scanner.current - scanner.start == start + length &&
      memcmp(scanner.start + start, rest, length) == 0) {
    return type;
  }

  return TOKEN_IDENTIFIER;
}

static TokenType identifierType() {
  switch (scanner.start[0]) {
    case '_': return checkKeyword(1, 0, "", TOKEN_UNDERSCORE);
    case 'a': 
      if (scanner.current - scanner.start > 1) {
        switch (scanner.start[1]) {
          case 'b': return checkKeyword(2, 6, "stract", TOKEN_ABSTRACT);
          case 's': return checkKeyword(2, 3, "ync", TOKEN_ASYNC);
        }
      }
      break;
    case 'b': return checkKeyword(1, 4, "reak", TOKEN_BREAK);
    case 'c':
      if (scanner.current - scanner.start > 1) {
        switch (scanner.start[1]) {
          case 'a': return checkKeyword(2, 3, "tch", TOKEN_CATCH);
          case 'l': return checkKeyword(2, 3, "ass", TOKEN_CLASS);
          case 'o': return checkKeyword(2, 6, "ntinue", TOKEN_CONTINUE);
        }
      }
      break;
    case 'd': return checkKeyword(1, 5, "elete", TOKEN_DELETE);
    case 'e': 
      if (scanner.current - scanner.start > 1) {
        switch (scanner.start[1]) {
          case 'l': return checkKeyword(2, 2, "se", TOKEN_ENUM);
          case 'n': return checkKeyword(2, 2, "um", TOKEN_ELSE);
        }
      }
      break;
    case 'f':
      if (scanner.current - scanner.start > 1) {
        switch (scanner.start[1]) {
          case 'a': return checkKeyword(2, 3, "lse", TOKEN_FALSE);
          case 'i': return checkKeyword(2, 5, "nally", TOKEN_FINALLY);
          case 'o': return checkKeyword(2, 1, "r", TOKEN_FOR);
          case 'u': return checkKeyword(2, 1, "n", TOKEN_FUN);
        }
      }
      break;
    case 'h': return checkKeyword(1, 3, "eap", TOKEN_HEAP);
    case 'i': 
      if (scanner.current - scanner.start > 1) {
        switch (scanner.start[1]) {
          case 'f': return checkKeyword(2, 0, "", TOKEN_IF);
          case 'n': return checkKeyword(2, 7, "terface", TOKEN_INTERFACE);
        }
      }
      break;
    case 'm': return checkKeyword(1, 3, "ark", TOKEN_MARK);
    case 'n': return checkKeyword(1, 3, "ull", TOKEN_NULL);
    case 'p': return checkKeyword(1, 4, "rint", TOKEN_PRINT);
    case 'r': return checkKeyword(1, 5, "eturn", TOKEN_RETURN);
    case 's': 
      if (scanner.current - scanner.start > 1) {
        switch (scanner.start[1]) {
          case 't': 
            if (scanner.current - scanner.start > 2) {
              switch (scanner.start[2]) {
                case 'a': 
                  if (scanner.current - scanner.start > 3) {
                    switch (scanner.start[3]) {
                      case 'c': return checkKeyword(4, 1, "k", TOKEN_STACK);
                      case 't': return checkKeyword(4, 2, "ic", TOKEN_STATIC);
                    }
                  }
                  break;
                case 't': return checkKeyword(3, 4, "ruct", TOKEN_STRUCT);
              }
            }
            break;
          case 'u': return checkKeyword(2, 3, "per", TOKEN_SUPER);
          case 'w': return checkKeyword(2, 4, "itch", TOKEN_SWITCH);
          case 'y': return checkKeyword(2, 2, "nc", TOKEN_SYNC);
        }
      }
      break;
    case 't':
      if (scanner.current - scanner.start > 1) {
        switch (scanner.start[1]) {
          case 'h': 
            if (scanner.current - scanner.start > 2) {
              switch (scanner.start[2]) {
                case 'i': return checkKeyword(3, 1, "s", TOKEN_THIS);
                case 'r': return checkKeyword(3, 2, "ow", TOKEN_THROW);
              }
            }
            break;
          case 'r': 
            if (scanner.current - scanner.start > 2) {
              switch (scanner.start[2]) {
                case 'u': return checkKeyword(3, 1, "e", TOKEN_TRUE);
                case 'y': return checkKeyword(3, 0, "", TOKEN_TRY);
              }
            }
            break;
        }
      }
      break;
    case 'v': return checkKeyword(1, 2, "ar", TOKEN_VAR);
    case 'w': return checkKeyword(1, 4, "hile", TOKEN_WHILE);
  }
  return TOKEN_IDENTIFIER;
}

static Token identifier() {
  while (isAlpha(peek()) || isDigit(peek() || peek() == '_')) advance();
  return makeToken(identifierType());
}

#pragma endregion

static Token number() {
  TokenType type = TOKEN_INT;
  while (isDigit(peek())) advance();

  // Look for a fractional part.
  if (peek() == '.' && isDigit(peekNext())) {
    type = TOKEN_FLT;
    advance(); // Consume the "."

    while (isDigit(peek())) advance();
    if(peek() == 'e') {
      if(isDigit(peekNext())) {
        advance(); // Consume the "e"
        while (isDigit(peek())) advance();
      } else if(peekNext() == '-' && isDigit(peekN(2))) {
        advanceN(2); // Consume the "e-"
        while (isDigit(peek())) advance();
      }
    }
  } else if(peek() == '`' && isDigit(peekNext())) {
    type = TOKEN_FIX;
    advance(); // Consume the "`"
    while (isDigit(peek())) advance();
  }

  return makeToken(type);
}

#pragma region String/Character

static int escape(int isString) {
    advance(); // Consume '\'
    switch(peek()) {
        case 'r': case 't': case 'n':
            advance();
            return true;
        case 'u':
            advance();
            return isHex(peek());
        case '\'':
            if(!isString) {
                advance();
                return true;
            }
            return false;
        case '"':
            if(isString) {
                advance();
                return true;
            }
            return false;
        default: return false;
    }
}

static Token string() {
 while (peek() != '"' && !isAtEnd()) {
    if (peek() == '\n') scanner.line++;
    else if(peek() == '\'' && !escape(true)) return errorToken("Invalid escape code.");
    else advance();
  }

  if (isAtEnd()) return errorToken("Unterminated string.");
  advance();
  return makeToken(TOKEN_STRING);
}

static Token character() {
  if(peek() == '\'') return errorToken("Empty character.");
  if(peek() == '\'' && !escape(false)) return errorToken("Invalid escape code.");
  else advance();

  if (isAtEnd() || peek() != '\'') return errorToken("Unterminated char.");
  advance();
  return makeToken(TOKEN_CHAR);
}

#pragma endregion

Token scanToken() {
  skipWhitespace();
  scanner.start = scanner.current;
  if (isAtEnd()) return makeToken(TOKEN_EOF);
  char c = advance();

  if (isAlpha(c) || c == '_') return identifier();
  else if (isDigit(c)) return number();
  
  switch (c) {
    case '(': return makeToken(TOKEN_LEFT_PAREN);
    case ')': return makeToken(TOKEN_RIGHT_PAREN);
    case '{': return makeToken(TOKEN_LEFT_BRACE);
    case '}': return makeToken(TOKEN_RIGHT_BRACE);
    case '[': return makeToken(TOKEN_LEFT_BRACKET);
    case ']': return makeToken(TOKEN_RIGHT_BRACKET);
    case ';': return makeToken(TOKEN_SEMICOLON);
    case ',': return makeToken(TOKEN_COMMA);
    case '.': return makeToken(TOKEN_DOT);
    case '$': return makeToken(TOKEN_DOLLAR);

    case '-': return makeToken(match('-') ? TOKEN_DEC : TOKEN_MINUS);
    case '+': return makeToken(match('+') ? TOKEN_INC : TOKEN_PLUS);
    case '/': return makeToken(TOKEN_SLASH);
    case '*': return makeToken(TOKEN_STAR);
    case '%': return makeToken(TOKEN_MODULO);

    case '&': return makeToken(match('&') ? TOKEN_AND : TOKEN_BIT_AND);
    case '|': return makeToken(match('|') ? TOKEN_OR : TOKEN_BIT_OR);
    case '^': return makeToken(TOKEN_BIT_XOR);

    case '!': return makeToken(match('=') ? TOKEN_BANG_EQUAL : TOKEN_BANG);
    case '=': return makeToken(match('=') ? TOKEN_EQUAL_EQUAL : TOKEN_EQUAL);
    case '<': return makeToken(match('=') ? (match('>') ? TOKEN_COMPARE : TOKEN_LESS_EQUAL) : (match('<') ? TOKEN_L_SHIFT : TOKEN_LESS));
    case '>': return makeToken(match('=') ? TOKEN_GREATER_EQUAL : (match('>') ? TOKEN_R_SHIFT : TOKEN_GREATER));
    case '@': return makeToken(TOKEN_AT);
    case '?': return makeToken(match('?') ? TOKEN_QMARK_QMARK : TOKEN_QMARK);
    case ':': return makeToken(match('?') ? TOKEN_COLON_COLON : TOKEN_COLON);
    case '#': return makeToken(TOKEN_HASH);

    case '"': return string();
    case '\'': return character();
  }

  return errorToken("Unexpected character.");
}
