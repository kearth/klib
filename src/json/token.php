<?php 

namespace Ktools\Json;

class Token {
    const T_BEGIN_DOC  = 0;
    const T_STRING     = 1;
    const T_NUMBER     = 2;
    const T_BOOLEAN    = 4;
    const T_NULL       = 8;
    const T_BEGIN_OBJ  = 16;
    const T_END_OBJ    = 32;
    const T_BEGIN_ARR  = 64;
    const T_END_ARR    = 128;
    const T_SEP_COLON  = 256;
    const T_SEP_COMMMA = 512;
    const T_SEP_QUOTE  = 1024;
    const T_END_DOC    = 2048;
}
