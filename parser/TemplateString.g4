
mode TemplateString;

fragment TS_PLAIN
    : ~['\\$]
    | '\\' ['\\bfnrtv$]
    | '\\' [uU] HEXDIGIT HEXDIGIT HEXDIGIT HEXDIGIT
    | '\\' [xX] HEXDIGIT HEXDIGIT
    ;

TS_RAW:         TS_PLAIN+;
TS_EXPR_START:  '${' -> pushMode(StrExpr);
TS_QUOTE:		QUOTE -> type(QUOTE), popMode;
TS_IDENTIFIER:  '$' IDENTIFIER;
