lexer grammar DefaultMode;

// 关键词
TRUE        : 'true';
FALSE       : 'false';
FOR         : 'for';
IN          : 'in';
IF          : 'if';
WHILE       : 'while';
DO          : 'do';
BREAK       : 'break';
CONTINUE    : 'continue';
FUNC        : 'func';
WHEN        : 'when';
ELSE        : 'else';
NIL         : 'nil';
UNDEFINED   : 'undefined';
RETURN_NONE : 'return' [ \t]* '\n';
RETURN      : 'return';
EXPORT      : 'export';
CLASS       : 'class';
DEFER       : 'defer';
BLOCK_DEFER : 'blockDefer';
THROW       : 'throw';
TRY         : 'try';
CATCH       : 'catch';
FINALLY     : 'finally';
STATIC      : 'static';
ASSERT      : 'assert';
EXTEND      : 'extend';
USE_AT      : 'use@';
USE         : 'use';
SWITCH      : 'switch';
CASE        : 'case';
FALLTHROUGH : 'fallthrough';
DEFAULT     : 'default';

// 数值

fragment DECDIGIT   : [0-9];
fragment HEXDIGIT   : [0-9a-zA-Z];
fragment OCTDIGIT   : [0-7];
fragment BINDIGIT   : [01];

WS          : [ \t\r\n]+    -> skip;
LINECOMMENT : '//' ~[\n]*   -> skip;
LINECOMMENT2: '#' ~[\n]*    -> skip;
BLOCKCOMMENT: '/*' .*? '*/' -> skip;
INT_ZERO    : '0';
INT_DEC     : [1-9] DECDIGIT*;
INT_HEX     : '0' [xX] HEXDIGIT+;
INT_OCT     : '0' OCTDIGIT+;
INT_BIN     : '0' [bB] BINDIGIT+;
BIGNUM      : ('0' | INT_DEC) ('.' DECDIGIT+)? ('L' | 'l');
FLOAT       : ('0' | INT_DEC) '.' DECDIGIT+;
fragment ESCCHAR
            : ['"\\bfnrtv]
            ;
fragment STRCHAR
            : ~['\\\r\n]
            | '\\' ESCCHAR
            | '\\' [uU] HEXDIGIT HEXDIGIT HEXDIGIT HEXDIGIT
            | '\\' [xX] HEXDIGIT HEXDIGIT
            ;
fragment RSTRCHAR
            : ~[']
            | '\\' '\''
            ;
STRING      : [rR] '\'' RSTRCHAR* '\'';

// 符号
MORE_ARGS   : '...';
LEAD_TO     : '->';
ARROW       : '=>';

// 运算符
POW                 : '**';
PLUS_PLUS           : '++';
MINUS_MINUS         : '--';
EQUAL               : '==';
NOT_EQUAL           : '!=';
GTEQ                : '>=';
LTEQ                : '<=';
LOCAL_ASSIGN        : ':=';
PLUS_ASSIGN         : '+=';
MINUS_ASSIGN        : '-=';
TIMES_ASSIGN        : '*=';
DIV_ASSIGN          : '/=';
MOD_ASSIGN          : '%=';
LOGIC_AND           : '&&';
LOGIC_OR            : '||';
OPTIONAL_CALL       : '?.';
OPTIONAL_ELSE       : '??';
BIT_AND             : '&';
BIT_OR              : '|';
BIT_NOT             : '~';
BIT_SHL             : '<<';
BIT_SHR             : '>>';
BIT_XOR             : '^';
BIT_AND_ASSIGN      : '&=';
BIT_OR_ASSIGN       : '|=';
BIT_SHL_ASSIGN      : '<<=';
BIT_SHR_ASSIGN      : '>>=';
BIT_XOR_ASSIGN      : '^=';
RANGE_WITHOUT_END   : '..<';
RANGE_WITH_END      : '..';

DOT             : '.';
COMMA           : ',';
SEMICOLON       : ';';
COLON           : ':';
L_PAREN         : '(';
R_PAREN         : ')';
L_CURLY         : '{';
R_CURLY         : '}';
L_BRACKET       : '[';
R_BRACKET       : ']';
LOGIC_NOT       : '!';
QUESTION        : '?';
GT              : '>';
LT              : '<';
ASSIGN          : '=';
PLUS            : '+';
MINUS           : '-';
TIMES           : '*';
DIV             : '/';
MOD             : '%';
SINGLE_AT       : '@';
DOUBLE_AT       : '@@';
QUOTE           : '\'' -> pushMode(TemplateString);

// DOLLAR          : '$';

fragment ID_STARTING: [a-zA-Z$_\u4E00-\u9FA5];
IDENTIFIER  : ID_STARTING (ID_STARTING|[0-9])*;
