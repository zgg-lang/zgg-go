import os.path as osp

outf = open(osp.join('parser', 'ZggLexer.g4'), 'w') 

# copy raws
default_lines = open(osp.join('parser', 'DefaultMode.g4')).readlines()
outf.write('lexer grammar ZggLexer;\n')
for l in default_lines[1:]:
    outf.write(l)
with open(osp.join('parser', 'TemplateString.g4')) as ts:
    outf.write(ts.read())

# read normal tokens
default_tokens = set()
max_len = 0
for l in default_lines[1:]:
    l = l.strip()
    if not l:
        continue
    start = l[0]
    if not 'A' <= start <= 'Z':
        continue
    token = l.split(':')[0].strip()
    if len(token) > max_len:
        max_len = len(token)
    default_tokens.add(token)

# write mode StrExpr
outf.write('\n\nmode StrExpr;\n\n')
tails = dict(
    L_CURLY='type(L_CURLY), pushMode(StrExpr)',
    R_CURLY='type(R_CURLY), popMode',
    QUOTE='type(QUOTE), pushMode(TemplateString)',
    WS='skip',
    LINE_COMMENT='skip',
    BLOCK_COMMENT=', skip',
)
for token in default_tokens:
    padding = ' ' * (max_len - len(token) + 1)
    outf.write('StrExpr_{0}{1}: {0}{1} -> '.format(token, padding))
    if token in tails:
        outf.write(tails[token])
    else:
        outf.write('type({0})'.format(token))
    outf.write(';\n')
