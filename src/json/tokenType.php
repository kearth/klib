<?php 

namespace Ktools\Json;

class TokenType {
    const T_STRING     = 1;
    const T_NUMBER     = 2;
    const T_BOOLEAN    = 4;
    const T_NULL       = 8;
    const T_ARRAY      = 16;
    const T_OBJECT     = 32;
    const T_END_DOC    = 64;
    const T_BEGIN_OBJ  = 128;
    const T_END_OBJ    = 256;
    const T_BEGIN_ARR  = 512;
    const T_END_ARR    = 1024;
    const T_SEP_COLON  = 2048;
    const T_SEP_COMMMA = 4096;
}
